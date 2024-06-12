package go_concurrent

import "context"

type Future[V any] interface {
	Get(ctx context.Context) (V, error)
	GetWithChan() <-chan V
	Cancel() bool
	Canceled() bool
	Done() bool
}
