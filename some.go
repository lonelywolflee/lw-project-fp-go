package lwfp

type Some[T any] struct {
	v T
}

func (s Some[T]) Map(fn func(T) any) (result Maybe[any]) {
	return Do(func() Maybe[any] {
		return Just(fn(s.Get().(T)))
	})
}

func (s Some[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any] {
	return Do(func() Maybe[any] {
		return fn(s.Get().(T))
	})
}

func (s Some[T]) Get() any {
	return s.v
}
