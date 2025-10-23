package maybe

// Maybe is a monad that represents an optional value or a computation that might fail.
// It provides a functional programming approach to handle nullable values and errors.
//
// The Maybe type has three concrete implementations:
//   - Some[T]: represents a value that exists
//   - None[T]: represents an absent value
//   - Failure[T]: represents an error state
//
// Design Note: Due to Go's type system constraints, methods Map and FlatMap operate on the same type T.
// For type conversions (T → R), use the helper functions: maybe.Map[T, R]() and maybe.FlatMap[T, R]()
//
// Example usage (same-type transformations):
//
//	result := Just(10).
//	    Map(func(x int) int { return x * 2 }).        // 10 → 20
//	    Filter(func(x int) bool { return x > 15 }).   // keeps 20
//	    Map(func(x int) int { return x + 5 })         // 20 → 25
//
// Example usage (type conversions with helper functions):
//
//	import "strconv"
//	result := Map(Just(42), strconv.Itoa)  // int → string: Just("42")
type Maybe[T any] interface {
	// Map transforms the value inside Maybe using the provided function.
	// The function must return the same type T.
	// For type conversion to a different type R, use the helper function: maybe.Map[T, R](m, fn)
	//
	// Behavior:
	//   - If Maybe is Some, applies the function and returns a new Some with the result
	//   - If Maybe is None, returns None (function not called)
	//   - If Maybe is Failure, returns Failure with same error (function not called)
	//   - If the function panics, it's caught and converted to a Failure
	//
	// Example:
	//
	//	result := Just(10).Map(func(x int) int { return x * 2 })  // Just(20)
	Map(fn func(T) T) Maybe[T]

	// MapIfEmpty provides both recovery and error transformation mechanisms for None states.
	// This method allows converting an empty Maybe into either a Some (recovery) or a Failure (error transformation).
	// The function returns (T, error) to support both recovery and error transformation patterns.
	//
	// Behavior:
	//   - If Maybe is Some: returns the original Some unchanged (function not called)
	//   - If Maybe is None: executes the function and wraps the result in a Maybe
	//   - If Maybe is Failure: returns the original Failure unchanged (function not called)
	//   - If the function returns (value, nil): returns Just(value) - recovery pattern
	//   - If the function returns (_, error): returns Failed(error) - error transformation pattern
	//   - If the function panics: catches the panic and returns Failure
	//
	// Use Cases:
	//   1. Recovery: Provide default values or fallback logic when a value is absent
	//   2. Error Transformation: Convert None to Failure with custom error
	//
	// Example:
	//
	//	// Use Case 1: Recovery - Provide a default value for None
	//	result := Empty[int]().MapIfEmpty(func() (int, error) {
	//	    return 42, nil
	//	}) // Just(42)
	//
	//	// Use Case 2: Error Transformation - Convert None to Failure
	//	result := Empty[int]().MapIfEmpty(func() (int, error) {
	//	    return 0, errors.New("value required")
	//	}) // Failed[int]("value required")
	//
	//	// Some remains unchanged
	//	result := Just(10).MapIfEmpty(func() (int, error) {
	//	    return 42, nil
	//	}) // Just(10) - function not called
	//
	//	// Failure remains unchanged
	//	result := Failed[int](err).MapIfEmpty(func() (int, error) {
	//	    return 42, nil
	//	}) // Failed[int](err) - function not called
	//
	//	// Recovery with potential failure
	//	result := Empty[string]().MapIfEmpty(func() (string, error) {
	//	    data, err := os.ReadFile("default.txt")
	//	    if err != nil {
	//	        return "", err
	//	    }
	//	    return string(data), nil
	//	}) // Just(fileContent) or Failed[string](err)
	MapIfEmpty(fn func() (T, error)) Maybe[T]

	// MapIfFailed provides both error recovery and error transformation mechanisms for Failure states.
	// This method allows converting a failed Maybe into either a Some (recovery) or a different Failure (error transformation).
	// The function receives the original error and returns (T, error) to support both recovery and transformation patterns.
	//
	// Behavior:
	//   - If Maybe is Some: returns the original Some unchanged (function not called)
	//   - If Maybe is None: returns the original None unchanged (function not called)
	//   - If Maybe is Failure: executes the function with the error and wraps the result in a Maybe
	//   - If the function returns (value, nil): returns Just(value) - recovery pattern
	//   - If the function returns (_, error): returns Failed(error) - error transformation pattern
	//   - If the function panics: catches the panic and returns Failure
	//
	// Use Cases:
	//   1. Error Recovery: Convert Failure to Some by providing fallback values or retry logic
	//   2. Error Transformation: Wrap, enrich, or convert errors (e.g., DB errors → domain errors)
	//
	// Example:
	//
	//	// Use Case 1: Recovery - Recover from specific errors
	//	result := Failed[int](errors.New("not found")).MapIfFailed(func(err error) (int, error) {
	//	    if errors.Is(err, ErrNotFound) {
	//	        return 0, nil  // Provide default value
	//	    }
	//	    return 0, err  // Propagate other errors
	//	}) // Just(0)
	//
	//	// Use Case 2: Error Transformation - Wrap or enrich errors
	//	result := Failed[int](dbErr).MapIfFailed(func(err error) (int, error) {
	//	    return 0, fmt.Errorf("domain error: %w", err)
	//	}) // Failed[int](wrapped error)
	//
	//	// Some remains unchanged
	//	result := Just(10).MapIfFailed(func(err error) (int, error) {
	//	    return 42, nil
	//	}) // Just(10) - function not called
	//
	//	// None remains unchanged
	//	result := Empty[int]().MapIfFailed(func(err error) (int, error) {
	//	    return 42, nil
	//	}) // Empty[int]() - function not called
	//
	//	// Recovery with retry logic
	//	result := fetchData().MapIfFailed(func(err error) (Data, error) {
	//	    log.Printf("First attempt failed: %v, retrying...", err)
	//	    return fetchDataFromCache()
	//	}) // Try cache if fetch fails
	MapIfFailed(fn func(error) (T, error)) Maybe[T]

	// FlatMap is similar to Map but expects the function to return a Maybe[T].
	// This prevents nested Maybe structures and is useful for chaining operations that might fail.
	// The function must return Maybe[T] (same type).
	// For type conversion to Maybe[R], use the helper function: maybe.FlatMap[T, R](m, fn)
	//
	// Behavior:
	//   - If Maybe is Some, applies the function and returns the resulting Maybe[T]
	//   - If Maybe is None, returns None (function not called)
	//   - If Maybe is Failure, returns Failure with same error (function not called)
	//   - If the function panics, it's caught and converted to a Failure
	//
	// Example:
	//
	//	result := Just(5).FlatMap(func(x int) Maybe[int] {
	//	    if x > 0 {
	//	        return Just(x * 2)
	//	    }
	//	    return Empty[int]()
	//	})  // Just(10)
	FlatMap(fn func(T) Maybe[T]) Maybe[T]

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
	//	result := Just(10).Then(func(x int) { fmt.Println(x) })     // prints 10, returns Just(10)
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
	//	value, err := Just(42).Get()             // value = 42, err = nil
	//	value, err := Empty[int]().Get()         // value = 0, err = nil
	//	value, err := Failed[int](someErr).Get() // value = 0, err = someErr
	Get() (T, error)

	// OrElseGet returns the value inside Maybe if it exists (Some case),
	// otherwise calls the provided function and returns its result (None or Failure case).
	// The function receives an error parameter: nil for None, actual error for Failure.
	// This allows error-aware default value computation.
	//
	// Example:
	//
	//	value := Just(42).OrElseGet(func(err error) int { return 0 })              // returns 42
	//	value := Empty[int]().OrElseGet(func(err error) int { return 0 })          // returns 0, err is nil
	//	value := Failed[int](err).OrElseGet(func(e error) int {
	//	    log.Printf("Error: %v", e)
	//	    return 0
	//	})  // returns 0, logs error
	OrElseGet(fn func(error) T) T

	// OrElseDefault returns the value inside Maybe if it exists (Some case),
	// otherwise returns the provided default value (None or Failure case).
	// This is useful for providing a constant default value when Maybe is empty or failed.
	//
	// Example:
	//
	//	value := Just(42).OrElseDefault(0)        // returns 42
	//	value := Empty[int]().OrElseDefault(0)    // returns 0
	//	value := Failed[int](err).OrElseDefault(0)  // returns 0
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
	//	result := Failed[int](err).MatchThen(
	//	    func(x int) { fmt.Printf("Value: %d\n", x) },
	//	    func() { fmt.Println("No value") },
	//	    func(err error) { fmt.Printf("Error: %v\n", err) },
	//	) // prints "Error: <error message>", returns Failed[int](err)
	MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T]
}
