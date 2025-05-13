//go:build !solution

package dupcall

import (
	"context"
	"slices"
	"sync"
)

type results struct {
	r   interface{}
	err error
}

type call struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	queue     []chan results
	m         sync.Mutex
}

type Call struct {
	token   sync.Mutex
	current *call
}

func (o *Call) Do(
	ctx context.Context,
	cb func(context.Context) (interface{}, error),
) (interface{}, error) {
	o.token.Lock()

	if o.current == nil || o.current.ctx.Err() != nil {
		ctxt, cancel := context.WithCancel(ctx)
		current := &call{ctx: ctxt, ctxCancel: cancel, queue: make([]chan results, 0), m: sync.Mutex{}}

		o.current = current

		// o.current.m.Lock()
		go func() {
			result, err := cb(ctxt)
			// fmt.Println("done")
			o.token.Lock()

			if o.current == nil {
				o.token.Unlock()
				return
			}

			o.current.m.Lock()
			for _, w := range o.current.queue {
				w <- results{result, err}
			}
			o.current.m.Unlock()

			o.current = nil

			o.token.Unlock()
		}()

	}
	now := o.current
	ch := make(chan results, 1)
	now.m.Lock()
	now.queue = append(now.queue, ch)
	now.m.Unlock()
	// fmt.Println("here")
	o.token.Unlock()

	select {
	case <-ctx.Done():
		now.m.Lock()
		for i, w := range now.queue {
			if w == ch {
				now.queue = slices.Delete(now.queue, i, i+1)
				break
			}
		}
		if len(now.queue) == 0 {
			now.ctxCancel()
		}
		now.m.Unlock()
		return nil, ctx.Err()
	case res := <-ch:
		return res.r, res.err
	}
}
