package go_concurrent

import (
	"testing"
	"time"
)

func TestNewCountDownLatch(t *testing.T) {
	count := uint32(3)
	cdl := NewCountDownLatch(count)

	for count > 0 {
		go func() {
			defer cdl.CountDown()
			t.Log("count down")
			time.Sleep(10 * time.Millisecond)
		}()
		count--
	}
	t.Log("Start Await")
	cdl.Await()
	t.Log("Finish Await")
}
