package queue

import (
	"context"
	"errors"
)

type MemoryQueue struct {
	ch     chan string   // Channel to hold queue items
	closed chan struct{} // Channel to signal closure
}

var ErrQueueClosed = errors.New("queue is closed")

func NewMemoryQueue(size int) *MemoryQueue {
	return &MemoryQueue{
		ch:     make(chan string, size),
		closed: make(chan struct{}),
	}
}

func (mq *MemoryQueue) Push(item string) error {
	select {
	case mq.ch <- item:
		return nil
	case <-mq.closed:
		return ErrQueueClosed
	}
}

func (mq *MemoryQueue) Pop(ctx context.Context) (string, error) {
	select {
	case item := <-mq.ch:
		return item, nil
	case <-ctx.Done(): // stateless cancellation, external context
		return "", ctx.Err()
	case <-mq.closed: // stateful closure
		return "", ErrQueueClosed
	}
}

func (mq *MemoryQueue) Len() int {
	return len(mq.ch)
}

func (mq *MemoryQueue) Close() {
	select {
	case <-mq.closed:
		// already closed, do nothing
	default:
		close(mq.closed)
	}
}
