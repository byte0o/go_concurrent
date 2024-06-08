package go_concurrent

import "sync/atomic"

type BarrierAction func()

type CyclicBarrier struct {
	parties, count uint32
	barrierCommand BarrierAction
	blockingQueue  chan struct{}
}

func NewCyclicBarrier(parties uint32, barrier BarrierAction) *CyclicBarrier {
	return &CyclicBarrier{
		parties:        parties,
		count:          parties,
		barrierCommand: barrier,
		blockingQueue:  make(chan struct{}),
	}
}

func (b *CyclicBarrier) Await() {
	for {
		current := atomic.LoadUint32(&b.parties)
		if current == 0 || current == 1 {
			goto RESET
		}
		newValue := current - 1
		if atomic.CompareAndSwapUint32(&b.parties, current, newValue) {
			<-b.blockingQueue
			return
		}
	}
RESET:
	defer b.breakBarrier()
	b.barrierCommand()
}

func (b *CyclicBarrier) Reset() {
	b.breakBarrier()
}

func (b *CyclicBarrier) breakBarrier() {
	b.count = b.parties
	bq := b.blockingQueue
	b.blockingQueue = make(chan struct{})
	close(bq)
}
