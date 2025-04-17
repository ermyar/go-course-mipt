//go:build !solution

package ratelimit

import (
	"context"
	"errors"
	"time"
)

// Limiter is precise rate limiter with context support.
type Limiter struct {
	lock   chan bool // aka mutex + data
	tokens chan int
	stop   chan struct{}
	last   chan time.Time // aka mutex + data
}

var ErrStopped = errors.New("limiter stopped")

// NewLimiter returns limiter that throttles rate of successful Acquire() calls
// to maxSize events at any given interval.
func NewLimiter(maxCount int, interval time.Duration) *Limiter {
	l := Limiter{lock: make(chan bool, 1), tokens: make(chan int, maxCount),
		stop: make(chan struct{}, 1), last: make(chan time.Time, 1)}
	for i := 1; i <= maxCount; i++ {
		l.tokens <- maxCount - i
	}
	l.lock <- false
	l.last <- time.Now()

	if interval == 0 {
		interval = time.Millisecond
	}

	go func() {
		mod := interval / time.Duration(maxCount)
		ticker := time.NewTicker(mod)
		defer ticker.Stop()
		for {
			select {
			case t := <-ticker.C:
				lst := <-l.last
				if t.Sub(lst) >= mod {
					select {
					case l.tokens <- 0:
					default:
					}
				}
				l.last <- lst
			case <-l.stop:
				return
			}
		}
	}()

	return &l
}

func (l *Limiter) Acquire(ctx context.Context) error {
	stopped := <-l.lock
	defer func() { l.lock <- stopped }()
	if stopped {
		return ErrStopped
	}

	select {
	case <-l.tokens:
		<-l.last
		l.last <- time.Now()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (l *Limiter) Stop() {
	<-l.lock
	defer func() { l.lock <- true }()

	l.stop <- struct{}{}

}
