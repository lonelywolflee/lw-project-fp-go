package lwfp

type Failure[T any] struct {
	e error
}

func (f Failure[T]) Map(fn func(T) any) Maybe[any] {
	return Fail[any](f.e)
}

func (f Failure[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any] {
	return Fail[any](f.e)
}

func (f Failure[T]) Get() any {
	return f.e
}
