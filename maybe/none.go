package maybe

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

// Filter ignores the given function and returns None.
// Since None has no value, there's nothing to filter.
//
// Example:
//
//	none := Empty[int]()
//	result := none.Filter(func(x int) bool { return x > 0 }) // Empty[int]()
func (n None[T]) Filter(fn func(T) bool) Maybe[T] {
	return n
}

// Then ignores the given function and returns None.
// Since None has no value, there's nothing to apply the function to.
//
// Example:
//
//	none := Empty[int]()
//	result := none.Then(func(x int) { println(x) }) // Empty[int]()
func (n None[T]) Then(fn func(T)) Maybe[T] {
	return n
}

// Get returns nil, indicating the absence of a value.
//
// Example:
//
//	none := Empty[int]()
//	value := none.Get() // returns nil
func (n None[T]) Get() (T, error) {
	var zero T
	return zero, nil
}

// OrElseGet calls the provided function and returns its result.
// Since None has no value, this method always executes the function to get a default value.
// The function receives nil as the error parameter, indicating "no error, just empty".
//
// Example:
//
//	none := Empty[int]()
//	result := none.OrElseGet(func(err error) int { return 10 }) // returns 10, err is nil
func (n None[T]) OrElseGet(fn func(error) T) T {
	return fn(nil)
}

// OrElseDefault returns the provided default value.
// Since None has no value, this method always returns the given default.
//
// Example:
//
//	none := Empty[int]()
//	result := none.OrElseDefault(10) // returns 10
func (n None[T]) OrElseDefault(v T) T {
	return v
}

// FailIfEmpty converts None to Failure with the provided error.
// This is useful for treating "empty" as an error condition in a processing pipeline.
//
// Example:
//
//	none := Empty[int]()
//	result := none.FailIfEmpty(errors.New("value required")) // returns Fail[int]("value required")
func (n None[T]) FailIfEmpty(err error) Maybe[T] {
	return Fail[T](err)
}

// MatchThen applies the given functions based on the type of Maybe.
// If Maybe is Some, the some function is called with the value inside Some.
// If Maybe is None, the none function is called.
// If Maybe is Failure, the failure function is called with the error inside Failure.
//
// Example:
//
//	none := Empty[int]()
//	result := none.MatchThen(func(x int) { println(x) }, func() { println("none") }, func(err error) { println(err) })                            // prints "none"
//	result := Fail[int](errors.New("failed")).MatchThen(func(x int) { println(x) }, func() { println("none") }, func(err error) { println(err) }) // prints "failed"
func (n None[T]) MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T] {
	return Do(func() Maybe[T] {
		noneFn()
		return n
	})
}
