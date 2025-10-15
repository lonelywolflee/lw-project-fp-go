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
