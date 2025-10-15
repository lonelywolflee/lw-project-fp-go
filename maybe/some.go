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
		return Just(fn(s.GetValue()))
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
		return fn(s.GetValue())
	})
}

// GetValue returns the value wrapped in Some.
// Unlike the Get() method from the Maybe interface, this returns the concrete type T
// without requiring type assertion.
//
// Example:
//
//	some := Just(42)
//	value := some.GetValue() // returns 42 as int (no type assertion needed)
func (s Some[T]) GetValue() T {
	return s.v
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
		if fn(s.GetValue()) {
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
		fn(s.GetValue())
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
	return s.GetValue(), nil
}

// OrElseGet returns the value inside Some.
// Since Some contains a value, the provided function is never called.
//
// Example:
//
//	some := Just(5)
//	result := some.OrElseGet(func() int { return 10 }) // returns 5 (not 10)
func (s Some[T]) OrElseGet(fn func() T) T {
	return s.GetValue()
}

// OrElseDefault returns the value inside Some.
// Since Some contains a value, the provided default value is ignored.
//
// Example:
//
//	some := Just(5)
//	result := some.OrElseDefault(10) // returns 5 (not 10)
func (s Some[T]) OrElseDefault(v T) T {
	return s.GetValue()
}

// MatchThen applies the given functions based on the type of Maybe.
// If Maybe is Some, the some function is called with the value inside Some.
// If Maybe is None, the none function is called.
// If Maybe is Failure, the failure function is called with the error inside Failure.
//
// Example:
//
//	some := Just(5)
//	result := some.MatchThen(func(x int) { println(x) }, func() { println("none") }, func(err error) { println(err) }) // prints 5
//	result := Empty[int]().MatchThen(func(x int) { println(x) }, func() { println("none") }, func(err error) { println(err) }) // prints "none"
//	result := Fail[int](errors.New("failed")).MatchThen(func(x int) { println(x) }, func() { println("none") }, func(err error) { println(err) }) // prints "failed"
func (s Some[T]) MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T] {
	return Do(func() Maybe[T] {
		someFn(s.GetValue())
		return s
	})
}
