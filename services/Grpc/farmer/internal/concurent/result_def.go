package concurent

type Result[T any] struct {
	Value T
	Error error
}
