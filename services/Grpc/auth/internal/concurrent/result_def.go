package concurrent

type Result[T any] struct {
	Value T
	Error error
}
