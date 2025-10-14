package lwfp

// Failure represents a Maybe that contains an error.
// It is one of the three concrete implementations of the Maybe interface.
// Failure wraps an error and propagates it through the computation chain.
// All operations on Failure preserve and propagate the error state,
// implementing the "railway-oriented programming" pattern for error handling.
type Failure[T any] struct {
	e error
}

// Map ignores the given function and propagates the error.
// Since Failure represents an error state, no transformation is applied.
// The error is preserved and wrapped in a new Failure.
//
// Example:
//
//	failure := Fail[int](errors.New("failed"))
//	result := failure.Map(func(x int) any { return x * 2 }) // Fail[any](error)
func (f Failure[T]) Map(fn func(T) any) Maybe[any] {
	return Fail[any](f.e)
}

// FlatMap ignores the given function and propagates the error.
// Since Failure represents an error state, no transformation is applied.
// The error is preserved and wrapped in a new Failure.
//
// Example:
//
//	failure := Fail[int](errors.New("failed"))
//	result := failure.FlatMap(func(x int) Maybe[any] {
//	    return Just(x * 2)
//	}) // Fail[any](error)
func (f Failure[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any] {
	return Fail[any](f.e)
}

// GetError returns the error wrapped in Failure.
// This method provides direct access to the error without type assertion.
//
// Example:
//
//	failure := Fail[int](errors.New("something went wrong"))
//	err := failure.GetError() // returns error directly (no type assertion needed)
func (f Failure[T]) GetError() error {
	return f.e
}

// Filter ignores the given function and returns Failure.
// Since Failure represents an error state, no filtering is applied.
// The error is preserved and wrapped in a new Failure.
//
// Example:
//
//	failure := Fail[int](errors.New("failed"))
//	result := failure.Filter(func(x int) bool { return x > 0 }) // Fail[int](error)
func (f Failure[T]) Filter(fn func(T) bool) Maybe[T] {
	return f
}

// Then ignores the given function and returns Failure.
// Since Failure represents an error state, no function application is performed.
// The error is preserved and wrapped in a new Failure.
//
// Example:
//
//	failure := Fail[int](errors.New("failed"))
//	result := failure.Then(func(x int) { println(x) }) // Fail[int](error)
func (f Failure[T]) Then(fn func(T)) Maybe[T] {
	return f
}
