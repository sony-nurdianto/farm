package concurent

import "context"

type Result[T any] struct {
	Value T
	Error error
}

func SendResult(
	ctx context.Context,
	send chan any,
	recv any,
) {
	select {
	case <-ctx.Done():
		return
	case send <- recv:
	}
}
