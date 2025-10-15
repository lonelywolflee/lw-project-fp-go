package maybe

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
// The error is preserved, and the type is kept as Failure[T].
//
// Example:
//
//	failure := Fail[int](errors.New("failed"))
//	result := failure.Map(func(x int) int { return x * 2 }) // Fail[int](error)
func (f Failure[T]) Map(fn func(T) T) Maybe[T] {
	return f
}

// FlatMap ignores the given function and propagates the error.
// Since Failure represents an error state, no transformation is applied.
// The error is preserved, and the type is kept as Failure[T].
//
// Example:
//
//	failure := Fail[int](errors.New("failed"))
//	result := failure.FlatMap(func(x int) Maybe[int] {
//	    return Just(x * 2)
//	}) // Fail[int](error)
func (f Failure[T]) FlatMap(fn func(T) Maybe[T]) Maybe[T] {
	return f
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

// Get returns the error wrapped in Failure.
// This method provides direct access to the error without type assertion.
//
// Example:
//
//	failure := Fail[int](errors.New("something went wrong"))
//	err := failure.Get() // returns error directly (no type assertion needed)
func (f Failure[T]) Get() (T, error) {
	var zero T
	return zero, f.e
}

// OrElseGet calls the provided function and returns its result.
// Since Failure represents an error state with no valid value, this method always executes the function to get a default value.
// The function receives the actual error, allowing error-aware default value computation.
//
// Example:
//
//	failure := Fail[int](errors.New("failed"))
//	result := failure.OrElseGet(func(err error) int {
//	    log.Printf("Error occurred: %v", err)
//	    return 10
//	}) // returns 10, logs the error
func (f Failure[T]) OrElseGet(fn func(error) T) T {
	return fn(f.e)
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

// FailIfEmpty returns the original Failure unchanged, preserving the original error.
// The provided function is not called because Failure already contains an error.
//
// Example:
//
//	failure := Fail[int](errors.New("database error"))
//	result := failure.FailIfEmpty(func() error { return errors.New("empty") }) // returns Fail[int]("database error"), function not called
func (f Failure[T]) FailIfEmpty(fn func() error) Maybe[T] {
	return f
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
