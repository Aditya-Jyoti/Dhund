package queue

import "context"

type Queue interface {
	Push(item string) error
	Pop(ctx context.Context) (error, string)
	Len() int
	Close()
}
