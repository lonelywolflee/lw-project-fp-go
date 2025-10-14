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

// OrElseGet calls the provided function and returns its result.
// Since Failure represents an error state with no valid value, this method always executes the function to get a default value.
//
// Example:
//
//	failure := Fail[int](errors.New("failed"))
//	result := failure.OrElseGet(func() int { return 10 }) // returns 10
func (f Failure[T]) OrElseGet(fn func() T) T {
	return fn()
}

// OrElseDefault returns the provided default value.
// Since Failure represents an error state with no valid value, this method always returns the given default.
//
// Example:
//
//	failure := Fail[int](errors.New("failed"))
//	result := failure.OrElseDefault(10) // returns 10
func (f Failure[T]) OrElseDefault(v T) T {
	return v
}

// MatchThen applies the given functions based on the type of Maybe.
// If Maybe is Some, the some function is called with the value inside Some.
// If Maybe is None, the none function is called.
// If Maybe is Failure, the failure function is called with the error inside Failure.
//
// Example:
//
//	failure := Fail[int](errors.New("failed"))
//	result := failure.MatchThen(func(x int) { println(x) }, func() { println("none") }, func(err error) { println(err) }) // prints "failed"
func (f Failure[T]) MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T] {
	return Do(func() Maybe[T] {
		failureFn(f.e)
		return f
	})
}
