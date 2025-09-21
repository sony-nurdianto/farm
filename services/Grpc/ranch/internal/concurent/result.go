package concurent

import "context"

type ResultType[T any] = result[T]

type Result[T any] interface {
	Err(err error) Result[T]
	Res(value T) Result[T]
	SendResult(ctx context.Context, send chan any)
}

type result[T any] struct {
	Value T
	Error error
}

func NewResult[T any]() Result[T] {
	return &result[T]{}
}

func (r *result[T]) Err(err error) Result[T] {
	r.Error = err
	return r
}

func (r *result[T]) Res(value T) Result[T] {
	r.Value = value
	return r
}

func (r result[T]) SendResult(
	ctx context.Context,
	send chan any,
) {
	select {
	case <-ctx.Done():
		return
	case send <- r:
	}
}
