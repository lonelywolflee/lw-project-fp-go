package maybe

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

	// Get returns the value and error from Maybe.
	// For Some: returns (value, nil)
	// For None: returns (zero value, nil)
	// For Failure: returns (zero value, error)
	// This provides a Go-idiomatic way to extract values with error handling.
	//
	// Example:
	//
	//	value, err := Just(42).Get()     // value = 42, err = nil
	//	value, err := Empty[int]().Get() // value = 0, err = nil
	//	value, err := Fail[int](someErr).Get() // value = 0, err = someErr
	Get() (T, error)

	// OrElseGet returns the value inside Maybe if it exists (Some case),
	// otherwise calls the provided function and returns its result (None or Failure case).
	// This is useful for providing a computed default value when Maybe is empty or failed.
	//
	// Example:
	//
	//	value := Just(42).OrElseGet(func() int { return 0 })  // returns 42
	//	value := Empty[int]().OrElseGet(func() int { return 0 })  // returns 0
	//	value := Fail[int](err).OrElseGet(func() int { return 0 })  // returns 0
	OrElseGet(fn func() T) T

	// OrElseDefault returns the value inside Maybe if it exists (Some case),
	// otherwise returns the provided default value (None or Failure case).
	// This is useful for providing a constant default value when Maybe is empty or failed.
	//
	// Example:
	//
	//	value := Just(42).OrElseDefault(0)  // returns 42
	//	value := Empty[int]().OrElseDefault(0)  // returns 0
	//	value := Fail[int](err).OrElseDefault(0)  // returns 0
	OrElseDefault(v T) T

	// MatchThen performs pattern matching on the Maybe type and executes the appropriate function for side effects.
	// This provides a type-safe way to handle all three Maybe states (Some, None, Failure) with custom behavior.
	// The function returns the original Maybe unchanged, making it suitable for chaining.
	// If any of the executed functions panics, the panic is caught and converted to a Failure.
	//
	// Parameters:
	//   - someFn: Function called when Maybe is Some, receives the wrapped value
	//   - noneFn: Function called when Maybe is None
	//   - failureFn: Function called when Maybe is Failure, receives the error
	//
	// Example:
	//
	//	result := Just(42).MatchThen(
	//	    func(x int) { fmt.Printf("Value: %d\n", x) },
	//	    func() { fmt.Println("No value") },
	//	    func(err error) { fmt.Printf("Error: %v\n", err) },
	//	) // prints "Value: 42", returns Just(42)
	//
	//	result := Empty[int]().MatchThen(
	//	    func(x int) { fmt.Printf("Value: %d\n", x) },
	//	    func() { fmt.Println("No value") },
	//	    func(err error) { fmt.Printf("Error: %v\n", err) },
	//	) // prints "No value", returns Empty[int]()
	//
	//	result := Fail[int](err).MatchThen(
	//	    func(x int) { fmt.Printf("Value: %d\n", x) },
	//	    func() { fmt.Println("No value") },
	//	    func(err error) { fmt.Printf("Error: %v\n", err) },
	//	) // prints "Error: <error message>", returns Fail[int](err)
	MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T]
}
