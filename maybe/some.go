package maybe

// Some represents a Maybe that contains a value.
// It is one of the three concrete implementations of the Maybe interface.
// Some wraps a non-nil value and provides transformation methods that operate on this value.
type Some[T any] struct {
	v T
}

// Map applies the given function to the value inside Some and wraps the result in a new Maybe.
// The function must return the same type T (for type conversion to other types, use the helper Map function).
// If the function panics, the panic is caught and converted to a Failure.
//
// Example:
//
//	some := Just(5)
//	result := some.Map(func(x int) int { return x * 2 }) // Just(10)
//
//	// For type conversion, use the helper function:
//	result := Map(Just(42), strconv.Itoa) // Just("42")
func (s Some[T]) Map(fn func(T) T) (result Maybe[T]) {
	return Do(func() Maybe[T] {
		return Just(fn(s.v))
	})
}

// MapIfEmpty returns the original Some unchanged since the value is present.
// The recovery function is not called because there is no empty state to recover from.
//
// Example:
//
//	some := Just(42)
//	result := some.MapIfEmpty(func() (int, error) {
//	    return 100, nil  // This function is never called
//	}) // Just(42)
func (s Some[T]) MapIfEmpty(fn func() (T, error)) Maybe[T] {
	return s
}

// MapIfFailed returns the original Some unchanged since there is no error state.
// The recovery function is not called because there is no failure to recover from.
//
// Example:
//
//	some := Just(42)
//	result := some.MapIfFailed(func(err error) (int, error) {
//	    return 100, nil  // This function is never called
//	}) // Just(42)
func (s Some[T]) MapIfFailed(fn func(error) (T, error)) Maybe[T] {
	return s
}

// FlatMap applies the given function to the value inside Some.
// Unlike Map, the function is expected to return a Maybe[T], which prevents nested Maybe structures.
// The function must return Maybe[T] (for type conversion, use the helper FlatMap function).
// If the function panics, the panic is caught and converted to a Failure.
//
// Example:
//
//	some := Just(5)
//	result := some.FlatMap(func(x int) Maybe[int] {
//	    if x > 0 {
//	        return Just(x * 2)
//	    }
//	    return Empty[int]()
//	})
//
//	// For type conversion, use the helper function:
//	result := FlatMap(Just(42), func(x int) Maybe[string] {
//	    return Just(strconv.Itoa(x))
//	}) // Just("42")
func (s Some[T]) FlatMap(fn func(T) Maybe[T]) Maybe[T] {
	return Do(func() Maybe[T] {
		return fn(s.v)
	})
}

// Filter applies the given function to the value inside Some and returns a new Maybe.
// If the function returns false, the value is discarded and the result is None.
// If the function panics, the panic is caught and converted to a Failure.
//
// Example:
//
//	some := Just(5)
//	result := some.Filter(func(x int) bool { return x > 0 }) // Just(5)
func (s Some[T]) Filter(fn func(T) bool) Maybe[T] {
	return Do(func() Maybe[T] {
		if fn(s.v) {
			return s
		}
		return Empty[T]()
	})
}

// Then applies the given function to the value inside Some.
// If the function panics, the panic is caught and converted to a Failure.
//
// Example:
//
//	some := Just(5)
//	result := some.Then(func(x int) { println(x) }) // Just(5)
func (s Some[T]) Then(fn func(T)) Maybe[T] {
	return Do(func() Maybe[T] {
		fn(s.v)
		return s
	})
}

// Get returns the value inside Some with presence flag true and no error.
//
// Example:
//
//	some := Just(5)
//	value, ok, err := some.Get() // returns 5, true, nil
func (s Some[T]) Get() (T, bool, error) {
	return s.v, true, nil
}

// OrElseGet returns the value inside Some.
// Since Some contains a value, the provided function is never called.
// The function parameter receives an error (nil for None, actual error for Failure),
// but for Some this function is not executed.
//
// Example:
//
//	some := Just(5)
//	result := some.OrElseGet(func(err error) int { return 10 }) // returns 5 (function not called)
func (s Some[T]) OrElseGet(fn func(error) T) T {
	return s.v
}

// OrElseDefault returns the value inside Some.
// Since Some contains a value, the provided default value is ignored.
//
// Example:
//
//	some := Just(5)
//	result := some.OrElseDefault(10) // returns 5 (not 10)
func (s Some[T]) OrElseDefault(v T) T {
	return s.v
}

// OrPanic returns the value inside Some.
// Since Some contains a value, this method never panics and simply returns the wrapped value.
//
// Example:
//
//	some := Just(42)
//	value := some.OrPanic() // returns 42, never panics
func (s Some[T]) OrPanic() T {
	return s.v
}

// MatchThen applies the given functions based on the type of Maybe.
// If Maybe is Some, the some function is called with the value inside Some.
// If Maybe is None, the none function is called.
// If Maybe is Failure, the failure function is called with the error inside Failure.
//
// Example:
//
//	some := Just(5)
//	result := some.MatchThen(func(x int) { println(x) }, func() { println("none") }, func(err error) { println(err) })                            // prints 5
//	result := Empty[int]().MatchThen(func(x int) { println(x) }, func() { println("none") }, func(err error) { println(err) })                    // prints "none"
//	result := Failed[int](errors.New("failed")).MatchThen(func(x int) { println(x) }, func() { println("none") }, func(err error) { println(err) }) // prints "failed"
func (s Some[T]) MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T] {
	return Do(func() Maybe[T] {
		someFn(s.v)
		return s
	})
}
