package go_concurrent

import (
	"sync/atomic"
)

/**
 * A counting semaphore.  Conceptually, a semaphore maintains a set of
 * permits.  Each {@link #acquire} blocks if necessary until a permit is
 * available, and then takes it.  Each {@link #release} adds a permit,
 * potentially releasing a blocking acquirer.
 * However, no actual permit objects are used; the {@code Semaphore} just
 * keeps a count of the number available and acts accordingly.
 *
 * <p>Semaphores are often used to restrict the number of goroutines than can
 * access some (physical or logical) resource. For example, here is
 * a class that uses a semaphore to control access to a pool of items:
 * <pre>
 */
type Semaphore struct {
	permits       int32
	blockingQueue chan struct{}
}

func NewSemaphore(permits uint32) *Semaphore {
	return &Semaphore{permits: int32(permits), blockingQueue: make(chan struct{})}
}

/**
 * Acquires a permit from this semaphore, blocking until one is
 * available.
 *
 * <p>Acquires a permit, if one is available and returns immediately,
 * reducing the number of available permits by one.
 *
 * <p>If no permit is available then the current goroutine becomes
 * disabled for goroutine scheduling purposes and lies dormant until
 * one of two things happens:
 * <ul>
 * <li>Some other goroutine invokes the {@link #release} method for this
 * semaphore and the current goroutine is next to be assigned a permit;
 * </ul>
 */
func (s *Semaphore) Acquire() {
	s.acquire(true)
}

/**
 * Acquires a permit from this semaphore, only if one is available at the
 * time of invocation.
 *
 * <p>Acquires a permit, if one is available and returns immediately,
 * with the value {@code true},
 * reducing the number of available permits by one.
 *
 * <p>If no permit is available then this method will return
 * immediately with the value {@code false}.
 *
 * <p>Even when this semaphore has been set to use a
 * fair ordering policy, a call to {@code TryAcquire()} <em>will</em>
 * immediately acquire a permit if one is available, whether or not
 * other goroutines are currently waiting.
 * This &quot;barging&quot; behavior can be useful in certain
 * circumstances, even though it breaks fairness. If you want to honor
 * the fairness setting
 */
func (s *Semaphore) TryAcquire() bool {
	return s.acquire(false)
}

func (s *Semaphore) acquire(blocking bool) bool {
	for {
		current := atomic.LoadInt32(&s.permits)
		newValue := current - 1
		if atomic.CompareAndSwapInt32(&s.permits, current, newValue) {
			if current == 0 || newValue < 0 {
				break
			}
			return true
		}
	}
	if !blocking {
		return false
	}
	<-s.blockingQueue
	return true
}

/**
 * Releases a permit, returning it to the semaphore.
 *
 * <p>Releases a permit, increasing the number of available permits by
 * one.  If any goroutines are trying to acquire a permit, then one is
 * selected and given the permit that was just released.  That goroutine
 * is (re)enabled for goroutine scheduling purposes.
 *
 * <p>There is no requirement that a goroutine that releases a permit must
 * have acquired that permit by calling {@link #acquire}.
 * Correct usage of a semaphore is established by programming convention
 * in the application.
 */
func (s *Semaphore) Release() {
	for {
		current := atomic.LoadInt32(&s.permits)
		newValue := current + 1
		if atomic.CompareAndSwapInt32(&s.permits, current, newValue) {
			if newValue <= 0 {
				s.blockingQueue <- struct{}{}
			}
			return
		}
	}
}
