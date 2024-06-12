package go_concurrent

import (
	"context"
	"errors"
	"sync/atomic"
)

var FutureTaskCancelledErr = errors.New("future task cancelled")

type FutureTask[V any] struct {
	v               chan V
	cancelled, done atomic.Bool
	err             error
	ctx             context.Context
	cancelFunc      context.CancelFunc
}

func newFutureTask[V any](ctx context.Context, cancelFunc context.CancelFunc) *FutureTask[V] {
	return &FutureTask[V]{v: make(chan V, 1), ctx: ctx, cancelFunc: cancelFunc}
}

func (ft *FutureTask[V]) Get() (V, error) {
	if ft.Cancelled() {
		var v V
		return v, FutureTaskCancelledErr
	}
	select {
	case v, _ := <-ft.v:
		return v, ft.err
	case <-ft.ctx.Done():
		var v V
		return v, ft.ctx.Err()
	}
}

func (ft *FutureTask[V]) GetWithChan() <-chan V {
	return ft.v
}

func (ft *FutureTask[V]) Done() bool {
	return ft.done.Load()
}
func (ft *FutureTask[V]) Cancel() bool {
	if ft.cancelled.CompareAndSwap(false, true) {
		ft.cancelFunc()
		return true
	}
	return false
}

func (ft *FutureTask[V]) Cancelled() bool {
	return ft.cancelled.Load()
}
