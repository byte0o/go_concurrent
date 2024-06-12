package go_concurrent

import (
	"context"
	"errors"
	"fmt"
)

type FutureGoTask[V any] func(ctx context.Context) (V, error)

func FutureGo[V any](ctx context.Context, cancelFunc context.CancelFunc, task FutureGoTask[V]) *FutureTask[V] {
	future := newFutureTask[V](cancelFunc)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				future.err = errors.New(fmt.Sprintf("%v", r))
			}
			future.done.Store(true)
		}()
		defer close(future.v)
		v, err := task(ctx)
		future.err = err
		if err == nil {
			future.v <- v
		}
	}()
	return future
}
