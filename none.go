package lwfp

type None[T any] struct {
}

func (n None[T]) Map(fn func(T) any) Maybe[any] {
	return Empty[any]()
}

func (n None[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any] {
	return Empty[any]()
}

func (n None[T]) Get() any {
	return nil
}
