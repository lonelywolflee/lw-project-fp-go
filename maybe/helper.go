package maybe

import (
	"errors"
	"fmt"
)

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
				result = Fail[T](err)
			} else {
				result = Fail[T](errors.New(fmt.Sprint(r)))
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
//	result := Map(Fail[int](err), strconv.Itoa) // Fail[string](err)
//
//	// Chaining with method calls
//	result := Map(
//	    Just(5).Filter(func(x int) bool { return x > 0 }),
//	    strconv.Itoa,
//	) // Just("5")
func Map[T any, R any](m Maybe[T], fn func(T) R) (output Maybe[R]) {
	value, err := m.Get()
	if err != nil {
		// Propagate Failure
		return Fail[R](err)
	}

	// Check if it's None by type assertion
	switch m.(type) {
	case None[T]:
		return Empty[R]()
	case Some[T]:
		// Apply the transformation with panic recovery
		return Do(func() Maybe[R] {
			return Just(fn(value))
		})
	default:
		// This should never happen, but return None as fallback
		return Empty[R]()
	}
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
//	        return Fail[int](err)
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
//	result := FlatMap(Fail[string](err), func(s string) Maybe[int] {
//	    return Just(len(s))
//	}) // Fail[int](err)
//
//	// Chaining for type conversion pipeline
//	result := FlatMap(
//	    Just("123"),
//	    func(s string) Maybe[int] {
//	        val, err := strconv.Atoi(s)
//	        if err != nil {
//	            return Fail[int](err)
//	        }
//	        return Just(val)
//	    },
//	) // Just(123)
func FlatMap[T any, R any](m Maybe[T], fn func(T) Maybe[R]) (output Maybe[R]) {
	value, err := m.Get()
	if err != nil {
		// Propagate Failure
		return Fail[R](err)
	}

	// Check if it's None by type assertion
	switch m.(type) {
	case None[T]:
		return Empty[R]()
	case Some[T]:
		// Apply the transformation with panic recovery
		return Do(func() Maybe[R] {
			return fn(value)
		})
	default:
		// This should never happen, but return None as fallback
		return Empty[R]()
	}
}
