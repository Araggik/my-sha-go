//go:build !solution

package rwmutex

// A RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// The zero value for a RWMutex is an unlocked mutex.
//
// If a goroutine holds a RWMutex for reading and another goroutine might
// call Lock, no goroutine should expect to be able to acquire a read lock
// until the initial read lock is released. In particular, this prohibits
// recursive read locking. This is to ensure that the lock eventually becomes
// available; a blocked Lock call excludes new readers from acquiring the
// lock.
type RWMutex struct {
	wCh    chan struct{}
	rCh    chan struct{}
	rCount int
}

// New creates *RWMutex.
func New() *RWMutex {
	mu := &RWMutex{
		wCh:    make(chan struct{}, 1),
		rCh:    make(chan struct{}, 1),
		rCount: 0,
	}

	mu.rCh <- struct{}{}

	return mu
}

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock
// call excludes new readers from acquiring the lock. See the
// documentation on the RWMutex type.
func (rw *RWMutex) RLock() {
	rw.wCh <- struct{}{}

	if rw.rCount == 0 {
		<-rw.rCh
	}

	rw.rCount++

	<-rw.wCh
}

// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading
// on entry to RUnlock.
func (rw *RWMutex) RUnlock() {
	rw.wCh <- struct{}{}

	if rw.rCount > 0 {
		rw.rCount--

		if rw.rCount == 0 {
			rw.rCh <- struct{}{}
		}
	}

	<-rw.wCh
}

// Lock locks rw for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (rw *RWMutex) Lock() {
	<-rw.rCh
}

// Unlock unlocks rw for writing. It is a run-time error if rw is
// not locked for writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular
// goroutine. One goroutine may RLock (Lock) a RWMutex and then
// arrange for another goroutine to RUnlock (Unlock) it.
func (rw *RWMutex) Unlock() {
	rw.rCh <- struct{}{}
}
