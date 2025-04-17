//go:build !solution

package tparallel

import (
	"sync"
)

type T struct {
	ch       chan int
	parallel bool
	token    sync.Mutex
	parent   *T
	wg       sync.WaitGroup
	localwg  sync.WaitGroup
}

func (t *T) Parallel() {
	t.parallel = true
	t.parent.wg.Done()
	t.localwg.Done()
	t.parent.token.Unlock()
	<-t.parent.ch
}

func (t *T) Run(subtest func(t *T)) {
	var subt T = T{
		ch:       make(chan int),
		token:    sync.Mutex{},
		parent:   t,
		parallel: false,
		wg:       sync.WaitGroup{},
		localwg:  sync.WaitGroup{},
	}

	defer func() {
		subt.localwg.Wait()
	}()

	subt.localwg.Add(1)
	t.wg.Add(1)
	t.localwg.Add(1)

	t.token.Lock()
	go func() {

		subtest(&subt)
		if !subt.parallel {
			t.token.Unlock()
			t.wg.Done()
			subt.localwg.Done()
		}
		close(subt.ch)
		subt.wg.Wait()
		t.localwg.Done()
	}()
}

func Run(topTests []func(t *T)) {
	var t T = T{
		ch:       make(chan int),
		parent:   nil,
		token:    sync.Mutex{},
		parallel: false,
		wg:       sync.WaitGroup{},
		localwg:  sync.WaitGroup{},
	}

	defer func() {
		t.wg.Wait()
		close(t.ch)
		t.localwg.Wait()
	}()

	for _, f := range topTests {
		t.token.Lock()
		subt := T{
			ch:       make(chan int),
			token:    sync.Mutex{},
			parent:   &t,
			parallel: false,
			wg:       sync.WaitGroup{},
			localwg:  sync.WaitGroup{},
		}
		t.wg.Add(1)
		t.localwg.Add(1)

		subt.localwg.Add(1)
		go func() {
			f(&subt)

			if !subt.parallel {
				t.token.Unlock()
				t.wg.Done()
				subt.localwg.Done()
			}
			close(subt.ch)

			subt.wg.Wait()
			t.localwg.Done()
		}()
		subt.localwg.Wait()
	}

}
