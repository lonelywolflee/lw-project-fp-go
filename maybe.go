package lwfp

type Maybe[T any] interface {
	Map(fn func(T) any) Maybe[any]
	FlatMap(fn func(T) Maybe[any]) Maybe[any]
	Get() any
}

func Just[T any](v T) Maybe[T] {
	return Some[T]{v: v}
}

func Empty[T any]() Maybe[T] {
	return None[T]{}
}

func Fail[T any](e error) Maybe[T] {
	return Failure[T]{e: e}
}
