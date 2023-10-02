package std

type Result[T any] struct {
	value T
	err   error
}

func (r Result[T]) IsErr() bool { return r.err != nil }
func (r Result[T]) Unwrap() T {
	if r.IsErr() {
		panic("IsErr()")
	}
	return r.value
}
func (r Result[T]) Error() error {
	if !r.IsErr() {
		panic("!IsErr()")
	}
	return r.err
}

func Ok[T any](value T) Result[T] {
	return Result[T]{value, nil}
}

func Err[T any](err error) Result[T] {
	return Result[T]{err: err}
}
