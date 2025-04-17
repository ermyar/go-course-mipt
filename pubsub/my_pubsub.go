//go:build !solution

package pubsub

import (
	"context"
	"fmt"
	"slices"
	"sync"
)

var _ Subscription = (*MySubscription)(nil)

type MySubscription struct {
	p       *MyPubSub
	subject string
	handler MsgHandler
}

type job struct {
	handler []*MySubscription
	msg     interface{}
}

func (s *MySubscription) Unsubscribe() {
	value := s.p.mp[s.subject]

	for i, ms := range value {
		if ms == s {
			value = slices.Delete(value, i, i+1)
		}
	}

	s.p.mp[s.subject] = value
}

var _ PubSub = (*MyPubSub)(nil)

type MyPubSub struct {
	mp       map[string][]*MySubscription
	meta     map[*MySubscription][]interface{}
	queue    chan *MySubscription
	token    sync.Mutex
	finished chan struct{}
}

func NewPubSub() PubSub {
	p := &MyPubSub{mp: make(map[string][]*MySubscription), queue: make(chan *MySubscription),
		finished: make(chan struct{}, 1), meta: make(map[*MySubscription][]interface{})}

	go func() {
		for {
			sub, exist := <-p.queue
			if exist {

				go func() {
					for {
						p.token.Lock()
						slice := p.meta[sub]
						p.token.Unlock()
						for _, msg := range slice {
							sub.handler(msg)
						}
						p.token.Lock()
						slice2 := p.meta[sub]
						slice2 = slice2[len(slice):]
						p.meta[sub] = slice2
						p.token.Unlock()
						if len(slice2) == 0 {
							break
						}
					}
				}()
			} else {
				return
			}
		}
	}()

	return p
}

func (p *MyPubSub) Subscribe(subj string, cb MsgHandler) (Subscription, error) {
	select {
	case <-p.finished:
		p.finished <- struct{}{}
		return nil, fmt.Errorf("already finished")
	default:
	}

	val := p.mp[subj]
	subs := &MySubscription{p: p, subject: subj, handler: cb}
	p.mp[subj] = append(val, subs)
	p.meta[subs] = nil

	return subs, nil
}

func (p *MyPubSub) Publish(subj string, msg interface{}) error {
	select {
	case <-p.finished:
		p.finished <- struct{}{}
		return fmt.Errorf("already finished")
	default:
	}

	val := p.mp[subj]

	for _, s := range val {
		p.token.Lock()
		tmp := p.meta[s]
		if len(tmp) == 0 {
			p.queue <- s
		}
		tmp = append(tmp, msg)
		p.meta[s] = tmp
		p.token.Unlock()
	}

	return nil
}

func (p *MyPubSub) Close(ctx context.Context) error {
	select {
	case <-p.finished:
		p.finished <- struct{}{}
		return fmt.Errorf("already finished")
	default:
	}

	defer func() {
		p.finished <- struct{}{}
		close(p.queue)
	}()

	if ctx.Done() == nil {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
}
