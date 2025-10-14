package lwfp

// None represents a Maybe that contains no value.
// It is one of the three concrete implementations of the Maybe interface.
// None represents the absence of a value without indicating an error.
// All operations on None return None, preserving the "empty" state.
type None[T any] struct {
}

// Map ignores the given function and returns Empty.
// Since None has no value, there's nothing to transform.
//
// Example:
//
//	none := Empty[int]()
//	result := none.Map(func(x int) any { return x * 2 }) // Empty[any]()
func (n None[T]) Map(fn func(T) any) Maybe[any] {
	return Empty[any]()
}

// FlatMap ignores the given function and returns Empty.
// Since None has no value, there's nothing to transform.
//
// Example:
//
//	none := Empty[int]()
//	result := none.FlatMap(func(x int) Maybe[any] {
//	    return Just(x * 2)
//	}) // Empty[any]()
func (n None[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any] {
	return Empty[any]()
}

// Get returns nil, indicating the absence of a value.
//
// Example:
//
//	none := Empty[int]()
//	value := none.Get() // returns nil
func (n None[T]) Get() any {
	return nil
}
