package otelog

import (
	"log"
	"sync/atomic"
	"testing"
	"time"
)

func getOrTimeout[T any](ch <-chan T, timeout time.Duration, t *testing.T) (res T, ok bool) {
	select {
	case res, ok = <-ch:
		return
	case <-time.After(timeout):
		t.Errorf("channel receive timeouted after %s", timeout)
		return
	}
}

func TestBatchLogProcessor_sendsAfterTimeout(t *testing.T) {
	processor := newBatchProcessor(WithBatchTimeout(5 * time.Millisecond))
	called := make(chan struct{})

	log.Printf("coucou")
	processor.batch(nil, func(batch []*LogRecord) {
		defer close(called)
		if len(batch) != 1 {
			t.Errorf("len(batch) = %d, wants 1", len(batch))
			return
		}
		if batch[0] != nil {
			t.Errorf("expected record to be nil")
		}
	})

	getOrTimeout(called, 10*time.Millisecond, t)

}

func TestBatchLogProcessor_sendsAftermaxQueueSize(t *testing.T) {
	processor := newBatchProcessor(WithMaxQueueSize(10))
	called := make(chan struct{})

	nbCalls := atomic.Int32{}
	callback := func(batch []*LogRecord) {
		defer close(called)

		calls := nbCalls.Add(1)
		if calls > 1 {
			t.Errorf("expected to be called once, called %d", calls)
		}

		if len(batch) != 10 {
			t.Errorf("len(batch) = %d, wants 10", len(batch))
			return
		}
		for i, r := range batch {
			if r != nil {
				t.Errorf("expected batch[[%d] to be nil", i)
			}
		}
	}

	for i := 0; i < 10; i++ {
		processor.batch(nil, callback)
	}

	getOrTimeout(called, 10*time.Millisecond, t)
	time.Sleep(10 * time.Millisecond)
}
