package go_concurrent

import "sync/atomic"

/**
 * A synchronization aid that allows one or more goroutines to wait until
 * a set of operations being performed in other goroutines completes.
 *
 * <p>A {@code CountDownLatch} is initialized with a given <em>count</em>.
 * The {@link #await await} methods block until the current count reaches
 * zero due to invocations of the {@link #countDown} method, after which
 * all waiting goroutines are released and any subsequent invocations of
 * {@link #await await} return immediately.  This is a one-shot phenomenon
 * -- the count cannot be reset.  If you need a version that resets the
 * count, consider using a {@link CyclicBarrier}.
 *
 * <p>A {@code CountDownLatch} is a versatile synchronization tool
 * and can be used for a number of purposes.  A
 * {@code CountDownLatch} initialized with a count of one serves as a
 * simple on/off latch, or gate: all goroutines invoking {@link #await await}
 * wait at the gate until it is opened by a goroutine invoking {@link
 * #countDown}.  A {@code CountDownLatch} initialized to <em>N</em>
 * can be used to make one goroutine wait until <em>N</em> goroutines have
 * completed some action, or some action has been completed N times.
 *
 * <p>A useful property of a {@code CountDownLatch} is that it
 * doesn't require that goroutines calling {@code countDown} wait for
 * the count to reach zero before proceeding, it simply prevents any
 * goroutine from proceeding past an {@link #await await} until all
 * goroutines could pass.
 *
 * <p><b>Sample usage:</b> Here is a pair of classes in which a group
 * of worker threads use two countdown latches:
 * <ul>
 * <li>The first is a start signal that prevents any worker from proceeding
 * until the driver is ready for them to proceed;
 * <li>The second is a completion signal that allows the driver to wait
 * until all workers have completed.
 * </ul>
 */
type CountDownLatch struct {
	count, max    uint32
	blockingQueue chan struct{}
}

func NewCountDownLatch(count uint32) *CountDownLatch {
	return &CountDownLatch{
		max:           count,
		blockingQueue: make(chan struct{}),
	}
}

/**
 * Decrements the count of the latch, releasing all waiting goroutines if
 * the count reaches zero.
 *
 * <p>If the current count is greater than zero then it is decremented.
 * If the new count is zero then all waiting goroutines are re-enabled for
 * goroutine scheduling purposes.
 *
 * <p>If the current count equals zero then nothing happens.
 */
func (c *CountDownLatch) CountDown() {
	count := atomic.AddUint32(&c.count, 1)
	if count > c.max {
		panic("too many calls to CountDown")
	}
	if count == c.max {
		close(c.blockingQueue)
	}
}

/*
*
  - Causes the current goroutine to wait until the latch has counted down to
  - zero, unless the goroutine is {@linkplain goroutine#interrupt interrupted}.
    *
  - <p>If the current count is zero then this method returns immediately.
    *
  - <p>If the current count is greater than zero then the current
  - goroutine becomes disabled for goroutine scheduling purposes and lies
  - dormant until one of two things happen:
  - <ul>
  - <li>The count reaches zero due to invocations of the
  - {@link #countDown} method;
  - </ul>
    *
*/
func (c *CountDownLatch) Await() {
	<-c.blockingQueue
}
