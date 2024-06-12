package go_concurrent

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFutureGo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	future := FutureGo[string](ctx, cancel, func(ctx context.Context) (string, error) {
		time.Sleep(200 * time.Millisecond)
		return "test", nil
	})

	for {
		if !future.Done() {
			time.Sleep(20 * time.Millisecond)
			t.Log("not done")
			continue
		}
		v, err := future.Get()
		assert.NoError(t, err)
		assert.Equal(t, "test", v)
		break
	}
}

func TestFutureGoWithCancel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	future := FutureGo[string](ctx, cancel, func(ctx context.Context) (string, error) {
		time.Sleep(200 * time.Millisecond)
		return "test", nil
	})

	i := 0
	for {
		i++
		if future.Cancelled() {
			t.Log("cancelled")
			t.Log(future.Get())

			return
		}
		if !future.Done() {
			time.Sleep(20 * time.Millisecond)
			t.Log("not done")
			if i == 3 {
				future.Cancel()
			}
			continue
		}

		v, err := future.Get()
		assert.NoError(t, err)
		assert.Equal(t, "test", v)
		break
	}
}
