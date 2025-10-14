package lwfp

// Maybe is a monad that represents an optional value or a computation that might fail.
// It provides a functional programming approach to handle nullable values and errors.
//
// The Maybe type has three concrete implementations:
//   - Some[T]: represents a value that exists
//   - None[T]: represents an absent value
//   - Failure[T]: represents an error state
//
// Example usage:
//
//	result := Just(10).
//	    Map(func(x int) any { return x * 2 }).
//	    FlatMap(func(x int) Maybe[any] { return Just(x + 5) })
type Maybe[T any] interface {
	// Map transforms the value inside Maybe using the provided function.
	// If Maybe is None or Failure, the transformation is skipped and the state is preserved.
	// If the function panics, it's caught and converted to a Failure.
	Map(fn func(T) any) Maybe[any]

	// FlatMap is similar to Map but expects the function to return a Maybe.
	// This is useful for chaining operations that might fail.
	// If Maybe is None or Failure, the transformation is skipped and the state is preserved.
	// If the function panics, it's caught and converted to a Failure.
	FlatMap(fn func(T) Maybe[any]) Maybe[any]

	// Filter applies a predicate function to the value inside Maybe.
	// If the predicate returns true, the Maybe is unchanged.
	// If the predicate returns false, the Maybe becomes None.
	// If Maybe is None or Failure, the predicate is not applied and the state is preserved.
	// If the function panics, it's caught and converted to a Failure.
	//
	// Example:
	//
	//	result := Just(10).Filter(func(x int) bool { return x > 5 }) // Just(10)
	//	result := Just(3).Filter(func(x int) bool { return x > 5 })  // Empty[int]()
	Filter(fn func(T) bool) Maybe[T]

	// Then applies a side-effect function to the value inside Maybe and returns the same Maybe.
	// This is useful for performing actions like logging or debugging without changing the value.
	// If Maybe is None or Failure, the function is not applied and the state is preserved.
	// If the function panics, it's caught and converted to a Failure.
	//
	// Example:
	//
	//	result := Just(10).Then(func(x int) { fmt.Println(x) }) // prints 10, returns Just(10)
	//	result := Empty[int]().Then(func(x int) { fmt.Println(x) }) // Empty[int](), nothing printed
	Then(fn func(T)) Maybe[T]
}

// Just creates a Maybe that contains a value (Some).
// Use this when you have a valid value to wrap.
//
// Example:
//
//	maybe := Just(42)
//	// To get the value with type safety, cast to Some and use GetValue():
//	value := maybe.(Some[int]).GetValue() // returns 42 as int
func Just[T any](v T) Maybe[T] {
	return Some[T]{v: v}
}

// Empty creates a Maybe that represents the absence of a value (None).
// Use this when you want to explicitly represent "no value" without an error.
//
// Example:
//
//	maybe := Empty[int]()
//	// Check if it's None:
//	if _, ok := maybe.(None[int]); ok {
//	    // Handle empty case
//	}
func Empty[T any]() Maybe[T] {
	return None[T]{}
}

// Fail creates a Maybe that represents an error state (Failure).
// Use this when you want to wrap an error in the Maybe monad.
//
// Example:
//
//	maybe := Fail[int](errors.New("something went wrong"))
//	// To get the error, cast to Failure and use GetError():
//	err := maybe.(Failure[int]).GetError() // returns the error
func Fail[T any](e error) Maybe[T] {
	return Failure[T]{e: e}
}
