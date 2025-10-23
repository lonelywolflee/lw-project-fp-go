package maybe

import (
	"errors"
	"fmt"
)

// ToMaybe converts Go's standard (value, error) tuple pattern to Maybe[T].
// This function bridges the gap between traditional Go error handling and the Maybe monad,
// making it easy to integrate existing Go APIs with functional programming patterns.
//
// The function takes a value and an error (the standard Go return pattern) and:
//   - Returns Failure[T] if error is not nil
//   - Returns Some[T] with the value if error is nil
//
// This is particularly useful when working with Go standard library functions
// or any function that returns (T, error).
//
// Example:
//
//	// Convert strconv.Atoi result directly to Maybe
//	result := ToMaybe(strconv.Atoi("42")) // Just(42)
//	result := ToMaybe(strconv.Atoi("abc")) // Fail(error)
//
//	// Use with file operations
//	content := ToMaybe(os.ReadFile("config.json"))
//
//	// Chain with other Maybe operations
//	parsed := ToMaybe(strconv.Atoi("123")).
//	    Filter(func(x int) bool { return x > 0 }).
//	    Map(func(x int) int { return x * 2 })
//
// Common use cases:
//   - Parsing functions (strconv.Atoi, json.Unmarshal, etc.)
//   - File I/O operations (os.ReadFile, os.Open, etc.)
//   - Network operations (http.Get, net.Dial, etc.)
//   - Database queries that return (result, error)
func ToMaybe[T any](v T, err error) Maybe[T] {
	if err != nil {
		return Failed[T](err)
	}
	return Just(v)
}

// Try executes a function that returns (T, error) and converts the result to Maybe[T].
// This function combines ToMaybe with panic recovery (via Do), providing both
// error handling and panic safety in a single operation.
//
// The function takes a function that returns (T, error) and:
//   - Returns Failure[T] if the function returns an error
//   - Returns Failure[T] if the function panics
//   - Returns Some[T] with the value if the function succeeds
//
// This is the most convenient way to wrap risky Go operations that might both
// return errors and potentially panic, converting them into the Maybe monad.
//
// Example:
//
//	// Safely parse with both error and panic protection
//	result := Try(func() (int, error) {
//	    return strconv.Atoi("42")
//	}) // Just(42)
//
//	// Handles errors
//	result := Try(func() (int, error) {
//	    return strconv.Atoi("not-a-number")
//	}) // Fail(error)
//
//	// Also catches panics
//	result := Try(func() (int, error) {
//	    var arr []int
//	    return arr[10], nil // This would panic, but Try catches it
//	}) // Fail(runtime error)
//
//	// Real-world example: JSON parsing with validation
//	config := Try(func() (Config, error) {
//	    data, err := os.ReadFile("config.json")
//	    if err != nil {
//	        return Config{}, err
//	    }
//	    var cfg Config
//	    err = json.Unmarshal(data, &cfg)
//	    return cfg, err
//	}).Filter(func(c Config) bool {
//	    return c.IsValid()
//	}).FailIfEmpty(func() error {
//	    return errors.New("invalid config")
//	})
//
// Difference between Try and ToMaybe:
//   - ToMaybe: Converts an already-computed (T, error) to Maybe[T]
//   - Try: Executes a function, handles both errors AND panics, returns Maybe[T]
//
// Use Try when:
//   - The operation might panic (e.g., slice access, type assertions)
//   - You want deferred execution with automatic error handling
//   - You need panic recovery in addition to error handling
//
// Use ToMaybe when:
//   - You already have a (T, error) tuple
//   - The operation is guaranteed not to panic
//   - You want more explicit control over execution
func Try[T any](fn func() (T, error)) Maybe[T] {
	return Do(func() Maybe[T] {
		return ToMaybe(fn())
	})
}

// Do executes the given function and catches any panics, converting them to Failure.
// This is a utility function that provides panic safety for operations that might fail.
// If the function panics with an error, that error is wrapped in a Failure.
// If the function panics with any other value, it's converted to an error and wrapped in a Failure.
//
// This function is used internally by Some.Map and Some.FlatMap to provide automatic
// error handling, but it can also be used directly for any risky operation.
//
// Example:
//
//	result := Do(func() Maybe[int] {
//	    // Some operation that might panic
//	    value := riskyOperation()
//	    return Just(value)
//	})
//	// If riskyOperation() panics, result will be a Failure containing the error
func Do[T any](fn func() Maybe[T]) (result Maybe[T]) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				result = Failed[T](err)
			} else {
				result = Failed[T](errors.New(fmt.Sprint(r)))
			}
		}
	}()

	return fn()
}

// Map transforms a Maybe[T] to Maybe[R] using the provided function.
// This is a helper function that enables type conversion across different types,
// which is not possible with the Maybe interface methods due to Go's type system constraints.
//
// The function takes a Maybe[T] and a transformation function, and returns Maybe[R].
// This allows converting values from one type to another while preserving the Maybe semantics.
//
// Error handling:
//   - If the input Maybe is None, returns None[R]
//   - If the input Maybe is Failure, returns Failure[R] with the same error
//   - If the input Maybe is Some but the function panics, returns Failure[R]
//   - If the input Maybe is Some and function succeeds, returns Some[R]
//
// Example:
//
//	// Convert int to string
//	result := Map(Just(42), strconv.Itoa) // Just("42")
//
//	// Convert with inline function
//	result := Map(Just(42), func(x int) string {
//	    return fmt.Sprintf("value: %d", x)
//	}) // Just("value: 42")
//
//	// Handles None - type is automatically converted
//	result := Map(Empty[int](), strconv.Itoa) // Empty[string]()
//
//	// Handles Failure - error is propagated with type conversion
//	result := Map(Failed[int](err), strconv.Itoa) // Failed[string](err)
//
//	// Chaining with method calls
//	result := Map(
//	    Just(5).Filter(func(x int) bool { return x > 0 }),
//	    strconv.Itoa,
//	) // Just("5")
func Map[T, R any](m Maybe[T], fn func(T) R) (output Maybe[R]) {
	m.MatchThen(
		func(v T) {
			output = Do(func() Maybe[R] {
				return Just(fn(v))
			})
		},
		func() {
			output = Empty[R]()
		},
		func(err error) {
			output = Failed[R](err)
		},
	)
	return
}

// FlatMap transforms a Maybe[T] to Maybe[R] using a function that returns Maybe[R].
// This is a helper function that enables type conversion with flatMapping across different types,
// which is not possible with the Maybe interface methods due to Go's type system constraints.
//
// Unlike Map, the transformation function returns Maybe[R] instead of R,
// which is useful for chaining operations that might fail or return empty.
// This prevents nested Maybe structures (Maybe[Maybe[R]]).
//
// The function takes a Maybe[T] and a transformation function that returns Maybe[R],
// and returns a flattened Maybe[R].
//
// Error handling:
//   - If the input Maybe is None, returns None[R]
//   - If the input Maybe is Failure, returns Failure[R] with the same error
//   - If the input Maybe is Some but the function panics, returns Failure[R]
//   - If the input Maybe is Some and function succeeds, returns the Maybe[R] from the function
//
// Example:
//
//	// Parse string to int with error handling
//	result := FlatMap(Just("42"), func(s string) Maybe[int] {
//	    val, err := strconv.Atoi(s)
//	    if err != nil {
//	        return Failed[int](err)
//	    }
//	    return Just(val)
//	}) // Just(42)
//
//	// Conditional transformation
//	result := FlatMap(Just("hello"), func(s string) Maybe[int] {
//	    if len(s) > 0 {
//	        return Just(len(s))
//	    }
//	    return Empty[int]()
//	}) // Just(5)
//
//	// Handles None - type is automatically converted
//	result := FlatMap(Empty[string](), func(s string) Maybe[int] {
//	    return Just(len(s))
//	}) // Empty[int]()
//
//	// Handles Failure - error is propagated with type conversion
//	result := FlatMap(Failed[string](err), func(s string) Maybe[int] {
//	    return Just(len(s))
//	}) // Failed[int](err)
//
//	// Chaining for type conversion pipeline
//	result := FlatMap(
//	    Just("123"),
//	    func(s string) Maybe[int] {
//	        val, err := strconv.Atoi(s)
//	        if err != nil {
//	            return Failed[int](err)
//	        }
//	        return Just(val)
//	    },
//	) // Just(123)
func FlatMap[T, R any](m Maybe[T], fn func(T) Maybe[R]) (output Maybe[R]) {
	m.MatchThen(
		func(v T) {
			output = Do(func() Maybe[R] {
				return fn(v)
			})
		},
		func() {
			output = Empty[R]()
		},
		func(err error) {
			output = Failed[R](err)
		},
	)
	return
}
