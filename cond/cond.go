//go:build !solution

package cond

// A Locker represents an object that can be locked and unlocked.
type Locker interface {
	Lock()
	Unlock()
}

// Cond implements a condition variable, a rendezvous point
// for goroutines waiting for or announcing the occurrence
// of an event.
//
// Each Cond has an associated Locker L (often a *sync.Mutex or *sync.RWMutex),
// which must be held when changing the condition and
// when calling the Wait method.
type Cond struct {
	L    Locker
	q    []chan int
	help chan int
}

// New returns a new Cond with Locker l.
func New(l Locker) *Cond {
	cond := Cond{L: l, q: make([]chan int, 0), help: make(chan int, 1)}
	cond.help <- 0
	return &cond
}

// Wait atomically unlocks c.L and suspends execution
// of the calling goroutine. After later resuming execution,
// Wait locks c.L before returning. Unlike in other systems,
// Wait cannot return unless awoken by Broadcast or Signal.
//
// Because c.L is not locked when Wait first resumes, the caller
// typically cannot assume that the condition is true when
// Wait returns. Instead, the caller should Wait in a loop:
//
//	c.L.Lock()
//	for !condition() {
//	    c.Wait()
//	}
//	... make use of condition ...
//	c.L.Unlock()
func (c *Cond) Wait() {
	c.L.Unlock()
	defer c.L.Lock()

	size := <-c.help
	ch := make(chan int, 1)
	c.q = append(c.q, ch)
	c.help <- size + 1
	<-ch
}

// Signal wakes one goroutine waiting on c, if there is any.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Signal() {
	size := <-c.help
	if size == 0 {
		c.help <- 0
		return
	}
	c.q[0] <- 0
	c.q = c.q[1:]
	c.help <- size - 1
}

// Broadcast wakes all goroutines waiting on c.
//
// It is allowed but not required for the caller to hold c.L
// during the call.
func (c *Cond) Broadcast() {
	for {
		size := <-c.help
		if size == 0 {
			c.help <- 0
			return
		}
		c.q[0] <- 0
		c.q = c.q[1:]
		c.help <- size - 1
	}
}
