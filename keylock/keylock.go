//go:build !solution

package keylock

import (
	"slices"
	"sync"
)

type job struct {
	str []string
	cnt int
	ch  chan struct{}
}

type KeyLock struct {
	mp    map[string][]*job
	token sync.Mutex
}

func New() *KeyLock {
	return &KeyLock{mp: make(map[string][]*job), token: sync.Mutex{}}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	j := job{str: keys, cnt: len(keys), ch: make(chan struct{}, 1)}

	l.token.Lock()
	// fmt.Println("hello2 ", keys)
	for _, s := range j.str {
		val, exist := l.mp[s]
		if !exist {
			val = make([]*job, 0)
		}
		val = append(val, &j)
		if len(val) == 1 {
			j.cnt--
			if j.cnt == 0 {
				j.ch <- struct{}{}
			}
		}
		l.mp[s] = val
	}
	l.token.Unlock()

	select {
	case <-j.ch:
		ff := func() {
			l.token.Lock()
			for _, s := range j.str {
				val := l.mp[s]
				val = val[1:]
				if len(val) > 0 {
					val[0].cnt--
					if val[0].cnt == 0 {
						val[0].ch <- struct{}{}
					}
				}
				l.mp[s] = val
				// fmt.Println("l.mp : ", s, val)
			}

			l.token.Unlock()
		}
		// fmt.Println(j.str)
		return false, ff
	case <-cancel:

		// fmt.Println("here", j.str)

		l.token.Lock()

		// fmt.Println("again?")

		for _, s := range j.str {
			val := l.mp[s]

			// fmt.Println("interrupt l.mp : ", s, val)

			if val[0] == &j {
				val = val[1:]
				if len(val) > 0 {
					val[0].cnt--
					if val[0].cnt == 0 {
						val[0].ch <- struct{}{}
					}
				}
			} else {
				for i := 1; i < len(val); i++ {
					// fmt.Printf("%p %p\n", val[i], &j)
					if val[i] == &j {
						val = slices.Delete(val, i, i+1)
						break
					}
				}
			}

			// fmt.Println("interrupt l.mp : ", s, val)
			l.mp[s] = val
		}

		l.token.Unlock()

		return true, nil
	}
}
