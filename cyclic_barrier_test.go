package go_concurrent

import (
	"testing"
	"time"
)

func TestNewCyclicBarrier(t *testing.T) {
	count := uint32(3)
	cb := NewCyclicBarrier(count, func() {
		t.Log("start exec  BarrierAction")
	})

	for count > 0 {
		go func(i uint32) {
			cb.Await()
			t.Log("start exec  goroutine", i)
		}(count)
		count--
	}
	time.Sleep(10 * time.Millisecond)
}
