package go_concurrent

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSemaphore_Acquire(t *testing.T) {
	s := NewSemaphore(3)

	s.Acquire()
	s.Acquire()
	s.Acquire()
	go func() {
		s.Release()
	}()
	s.Acquire()
}

func TestSemaphore_TryAcquire(t *testing.T) {
	s := NewSemaphore(3)
	s.Acquire()
	s.Acquire()
	assert.Equal(t, s.TryAcquire(), true)
	assert.Equal(t, s.TryAcquire(), false)
}
