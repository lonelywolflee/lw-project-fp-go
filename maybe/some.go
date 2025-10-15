package maybe

// Some represents a Maybe that contains a value.
// It is one of the three concrete implementations of the Maybe interface.
// Some wraps a non-nil value and provides transformation methods that operate on this value.
type Some[T any] struct {
	v T
}

// Map applies the given function to the value inside Some and wraps the result in a new Maybe.
// If the function panics, the panic is caught and converted to a Failure.
//
// Example:
//
//	some := Just(5)
//	result := some.Map(func(x int) any { return x * 2 }) // Just(10)
func (s Some[T]) Map(fn func(T) any) (result Maybe[any]) {
	return Do(func() Maybe[any] {
		return Just(fn(s.v))
	})
}

// FlatMap applies the given function to the value inside Some.
// Unlike Map, the function is expected to return a Maybe, which prevents nested Maybe structures.
// If the function panics, the panic is caught and converted to a Failure.
//
// Example:
//
//	some := Just(5)
//	result := some.FlatMap(func(x int) Maybe[any] {
//	    if x > 0 {
//	        return Just(x * 2)
//	    }
//	    return Empty[any]()
//	})
func (s Some[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any] {
	return Do(func() Maybe[any] {
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

// Get returns the value inside Some.
//
// Example:
//
//	some := Just(5)
//	value := some.Get() // returns 5
func (s Some[T]) Get() (T, error) {
	return s.v, nil
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

// FailIfEmpty returns the original Some unchanged since it contains a value.
// The provided error is ignored because Some is not empty.
//
// Example:
//
//	some := Just(5)
//	result := some.FailIfEmpty(errors.New("empty")) // returns Just(5), error ignored
func (s Some[T]) FailIfEmpty(err error) Maybe[T] {
	return s
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
//	result := Fail[int](errors.New("failed")).MatchThen(func(x int) { println(x) }, func() { println("none") }, func(err error) { println(err) }) // prints "failed"
func (s Some[T]) MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T] {
	return Do(func() Maybe[T] {
		someFn(s.v)
		return s
	})
}
