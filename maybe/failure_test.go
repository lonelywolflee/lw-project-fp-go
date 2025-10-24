package maybe_test

import (
	"errors"
	"testing"

	"github.com/lonelywolflee/lw-project-fp-go/maybe"
)

// TestFailure_GetError is removed - use TestFailure_Get instead
// Get() is now the unified interface method for accessing errors

func TestFailure_Get(t *testing.T) {
	t.Run("returns zero value and error for int", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		value, _, returnedErr := failure.Get()

		if returnedErr != err {
			t.Errorf("expected error %v, got %v", err, returnedErr)
		}
		if value != 0 {
			t.Errorf("expected zero value (0), got %d", value)
		}
	})

	t.Run("returns zero value and error for string", func(t *testing.T) {
		err := errors.New("database error")
		failure := maybe.Failed[string](err)
		value, _, returnedErr := failure.Get()

		if returnedErr != err {
			t.Errorf("expected error %v, got %v", err, returnedErr)
		}
		if value != "" {
			t.Errorf("expected zero value (empty string), got %s", value)
		}
	})

	t.Run("returns zero value and error for bool", func(t *testing.T) {
		err := errors.New("validation error")
		failure := maybe.Failed[bool](err)
		value, _, returnedErr := failure.Get()

		if returnedErr != err {
			t.Errorf("expected error %v, got %v", err, returnedErr)
		}
		if value != false {
			t.Errorf("expected zero value (false), got %v", value)
		}
	})

	t.Run("returns nil pointer and error for pointer type", func(t *testing.T) {
		err := errors.New("not found")
		failure := maybe.Failed[*int](err)
		value, _, returnedErr := failure.Get()

		if returnedErr != err {
			t.Errorf("expected error %v, got %v", err, returnedErr)
		}
		if value != nil {
			t.Errorf("expected nil pointer, got %v", value)
		}
	})

	t.Run("returns zero value struct and error", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		err := errors.New("user not found")
		failure := maybe.Failed[User](err)
		value, _, returnedErr := failure.Get()

		if returnedErr != err {
			t.Errorf("expected error %v, got %v", err, returnedErr)
		}
		if value.Name != "" || value.Age != 0 {
			t.Errorf("expected zero value User{}, got %+v", value)
		}
	})

	t.Run("preserves different error types", func(t *testing.T) {
		type CustomError struct {
			Code    int
			Message string
		}

		// Implement error interface
		errImpl := func(ce CustomError) error {
			return errors.New(ce.Message)
		}

		customErr := errImpl(CustomError{Code: 404, Message: "Not Found"})
		failure := maybe.Failed[int](customErr)
		_, _, returnedErr := failure.Get()

		if returnedErr != customErr {
			t.Errorf("expected custom error %v, got %v", customErr, returnedErr)
		}
	})
}

func TestFailure_Map(t *testing.T) {
	t.Run("propagates error and ignores function", func(t *testing.T) {
		err := errors.New("original error")
		failure := maybe.Failed[int](err)
		result := failure.Map(func(x int) int { return x * 2 })

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Failure.Map should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		executed := false
		result := failure.Map(func(x int) int {
			executed = true
			return x * 2
		})

		if executed {
			t.Error("Failure.Map should not execute the function")
		}

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Failure.Map should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		result := failure.Map(func(x int) int {
			panic("this should never be called")
		})

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Failure.Map should return Failure type without executing function")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("preserves error through multiple Map calls", func(t *testing.T) {
		err := errors.New("persistent error")
		result := maybe.Failed[int](err).
			Map(func(x int) int { return x * 2 })

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("chained Map should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})
}

func TestFailure_FlatMap(t *testing.T) {
	t.Run("propagates error and ignores function", func(t *testing.T) {
		err := errors.New("original error")
		failure := maybe.Failed[int](err)
		result := failure.FlatMap(func(x int) maybe.Maybe[int] {
			return maybe.Just(x * 2)
		})

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Failure.FlatMap should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		executed := false
		result := failure.FlatMap(func(x int) maybe.Maybe[int] {
			executed = true
			return maybe.Just(x * 2)
		})

		if executed {
			t.Error("Failure.FlatMap should not execute the function")
		}

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Failure.FlatMap should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[string](err)
		result := failure.FlatMap(func(x string) maybe.Maybe[string] {
			panic("this should never be called")
		})

		resultFailure, ok := result.(maybe.Failure[string])
		if !ok {
			t.Fatal("Failure.FlatMap should return Failure type without executing function")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("preserves error through multiple FlatMap calls", func(t *testing.T) {
		err := errors.New("persistent error")
		result := maybe.Failed[int](err).
			FlatMap(func(x int) maybe.Maybe[int] { return maybe.Just(x * 2) })

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("chained FlatMap should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("preserves error through mixed operations", func(t *testing.T) {
		err := errors.New("persistent error")
		result := maybe.Failed[int](err).
			Map(func(x int) int { return x * 2 })

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("mixed operations should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("can be used in railway-oriented programming", func(t *testing.T) {
		err := errors.New("validation failed")
		result := maybe.Failed[int](err).
			FlatMap(func(x int) maybe.Maybe[int] {
				if x > 0 {
					return maybe.Just(x * 2)
				}
				return maybe.Empty[int]()
			})

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("railway pattern should preserve Failure")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})
}

func TestFailure_Filter(t *testing.T) {
	t.Run("propagates error and ignores predicate", func(t *testing.T) {
		err := errors.New("original error")
		failure := maybe.Failed[int](err)
		result := failure.Filter(func(x int) bool { return x > 5 })

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Failure.Filter should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("does not execute the predicate function", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		executed := false
		result := failure.Filter(func(x int) bool {
			executed = true
			return true
		})

		if executed {
			t.Error("Failure.Filter should not execute the predicate function")
		}

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Failure.Filter should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("predicate can panic but never called", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[string](err)
		result := failure.Filter(func(x string) bool {
			panic("this should never be called")
		})

		resultFailure, ok := result.(maybe.Failure[string])
		if !ok {
			t.Fatal("Failure.Filter should return Failure type without executing predicate")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		err := errors.New("persistent error")
		result := maybe.Failed[int](err).
			Filter(func(x int) bool { return x > 5 }).
			Map(func(x int) int { return x * 2 })

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("chained Filter and Map should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("preserves error in railway pattern", func(t *testing.T) {
		err := errors.New("validation failed")
		result := maybe.Failed[int](err).
			Filter(func(x int) bool { return x > 0 })

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Filter should preserve Failure in railway pattern")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})
}

func TestFailure_Then(t *testing.T) {
	t.Run("propagates error and ignores function", func(t *testing.T) {
		err := errors.New("original error")
		failure := maybe.Failed[int](err)
		result := failure.Then(func(x int) { /* no-op */ })

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Failure.Then should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		executed := false
		result := failure.Then(func(x int) {
			executed = true
		})

		if executed {
			t.Error("Failure.Then should not execute the function")
		}

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Failure.Then should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[string](err)
		result := failure.Then(func(x string) {
			panic("this should never be called")
		})

		resultFailure, ok := result.(maybe.Failure[string])
		if !ok {
			t.Fatal("Failure.Then should return Failure type without executing function")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		err := errors.New("persistent error")
		var sideEffect int

		result := maybe.Failed[int](err).
			Then(func(x int) { sideEffect = x }).
			Map(func(x int) int { return x * 2 })

		if sideEffect != 0 {
			t.Errorf("side effect should not be triggered, got %d", sideEffect)
		}

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("chained Then and Map should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("preserves error through multiple Then calls", func(t *testing.T) {
		err := errors.New("validation failed")
		var callCount int

		result := maybe.Failed[int](err).
			Then(func(x int) { callCount++ }).
			Then(func(x int) { callCount++ }).
			Then(func(x int) { callCount++ })

		if callCount != 0 {
			t.Errorf("no Then calls should execute, got %d calls", callCount)
		}

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("multiple Then calls should preserve Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("preserves error in railway pattern with Then", func(t *testing.T) {
		err := errors.New("validation failed")
		result := maybe.Failed[int](err).
			Then(func(x int) { /* log */ })

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Then should preserve Failure in railway pattern")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})
}

func TestFailure_OrElseGet(t *testing.T) {
	t.Run("calls function and returns result", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		called := false
		var receivedErr error
		result := failure.OrElseGet(func(e error) int {
			called = true
			receivedErr = e
			return 42
		})

		if !called {
			t.Error("OrElseGet should call the function when Failure has no value")
		}
		if receivedErr != err {
			t.Errorf("expected error %v, got %v", err, receivedErr)
		}
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("returns string from function", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[string](err)
		result := failure.OrElseGet(func(e error) string { return "default" })

		if result != "default" {
			t.Errorf("expected 'default', got %s", result)
		}
	})

	t.Run("executes function every time", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		callCount := 0

		result1 := failure.OrElseGet(func(e error) int {
			callCount++
			return callCount
		})
		result2 := failure.OrElseGet(func(e error) int {
			callCount++
			return callCount
		})

		if result1 != 1 {
			t.Errorf("first call expected 1, got %d", result1)
		}
		if result2 != 2 {
			t.Errorf("second call expected 2, got %d", result2)
		}
		if callCount != 2 {
			t.Errorf("expected 2 function calls, got %d", callCount)
		}
	})

	t.Run("works with different types", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[[]int](err)
		result := failure.OrElseGet(func(e error) []int { return []int{1, 2, 3} })

		if len(result) != 3 {
			t.Errorf("expected slice length 3, got %d", len(result))
		}
		if result[0] != 1 || result[1] != 2 || result[2] != 3 {
			t.Errorf("expected [1 2 3], got %v", result)
		}
	})

	t.Run("can return zero values", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		result := failure.OrElseGet(func(e error) int { return 0 })

		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("function can compute complex values", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		result := failure.OrElseGet(func(e error) int {
			sum := 0
			for i := 1; i <= 10; i++ {
				sum += i
			}
			return sum
		})

		if result != 55 {
			t.Errorf("expected 55 (sum of 1-10), got %d", result)
		}
	})

	t.Run("different errors still call function", func(t *testing.T) {
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")

		failure1 := maybe.Failed[int](err1)
		failure2 := maybe.Failed[int](err2)

		result1 := failure1.OrElseGet(func(e error) int { return 10 })
		result2 := failure2.OrElseGet(func(e error) int { return 20 })

		if result1 != 10 {
			t.Errorf("expected 10, got %d", result1)
		}
		if result2 != 20 {
			t.Errorf("expected 20, got %d", result2)
		}
	})

	t.Run("function receives actual error", func(t *testing.T) {
		originalErr := errors.New("database connection failed")
		failure := maybe.Failed[int](originalErr)
		var capturedErr error
		result := failure.OrElseGet(func(err error) int {
			capturedErr = err
			return 100
		})

		if capturedErr != originalErr {
			t.Errorf("expected error %v, got %v", originalErr, capturedErr)
		}
		if result != 100 {
			t.Errorf("expected 100, got %d", result)
		}
	})

	t.Run("can use error in computation", func(t *testing.T) {
		err := errors.New("value must be positive")
		failure := maybe.Failed[int](err)
		result := failure.OrElseGet(func(e error) int {
			// Use error information to determine default
			if e.Error() == "value must be positive" {
				return 1 // Return positive default
			}
			return 0
		})

		if result != 1 {
			t.Errorf("expected 1, got %d", result)
		}
	})
}

func TestFailure_OrElseDefault(t *testing.T) {
	t.Run("returns the default value", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		result := failure.OrElseDefault(42)

		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("returns string default value", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[string](err)
		result := failure.OrElseDefault("default")

		if result != "default" {
			t.Errorf("expected 'default', got %s", result)
		}
	})

	t.Run("works with different types", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[[]int](err)
		result := failure.OrElseDefault([]int{1, 2, 3})

		if len(result) != 3 {
			t.Errorf("expected slice length 3, got %d", len(result))
		}
		if result[0] != 1 || result[1] != 2 || result[2] != 3 {
			t.Errorf("expected [1 2 3], got %v", result)
		}
	})

	t.Run("can return zero values", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		result := failure.OrElseDefault(0)

		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("returns same default every time", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		result1 := failure.OrElseDefault(10)
		result2 := failure.OrElseDefault(10)

		if result1 != 10 || result2 != 10 {
			t.Errorf("expected both results to be 10, got %d and %d", result1, result2)
		}
	})

	t.Run("different errors still return default", func(t *testing.T) {
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")

		failure1 := maybe.Failed[int](err1)
		failure2 := maybe.Failed[int](err2)

		result1 := failure1.OrElseDefault(100)
		result2 := failure2.OrElseDefault(100)

		if result1 != 100 || result2 != 100 {
			t.Errorf("expected both results to be 100, got %d and %d", result1, result2)
		}
	})
}

func TestFailure_MatchThen(t *testing.T) {
	t.Run("executes failureFn and returns original Failure", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		someCalled := false
		noneCalled := false
		failureCalled := false
		var capturedError error

		result := failure.MatchThen(
			func(x int) { someCalled = true },
			func() { noneCalled = true },
			func(e error) {
				failureCalled = true
				capturedError = e
			},
		)

		if someCalled {
			t.Error("someFn should not be called")
		}
		if noneCalled {
			t.Error("noneFn should not be called")
		}
		if !failureCalled {
			t.Error("failureFn should be called")
		}
		if capturedError != err {
			t.Errorf("expected captured error %v, got %v", err, capturedError)
		}

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MatchThen should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("can be used for error logging", func(t *testing.T) {
		err := errors.New("database error")
		failure := maybe.Failed[string](err)
		var log string

		result := failure.MatchThen(
			func(x string) { log = "Got value: " + x },
			func() { log = "No value" },
			func(e error) { log = "Error: " + e.Error() },
		)

		if log != "Error: database error" {
			t.Errorf("expected 'Error: database error', got %s", log)
		}

		resultFailure, ok := result.(maybe.Failure[string])
		if !ok {
			t.Fatal("MatchThen should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("catches panic in failureFn and converts to Failure", func(t *testing.T) {
		err := errors.New("original error")
		failure := maybe.Failed[int](err)

		result := failure.MatchThen(
			func(x int) {},
			func() {},
			func(e error) { panic("failureFn panic") },
		)

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MatchThen should return Failure when failureFn panics")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr.Error() != "failureFn panic" {
			t.Errorf("expected panic message, got %s", gotErr.Error())
		}
	})

	t.Run("someFn and noneFn can panic but never called", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		failureCalled := false

		result := failure.MatchThen(
			func(x int) { panic("someFn should not be called") },
			func() { panic("noneFn should not be called") },
			func(e error) { failureCalled = true },
		)

		if !failureCalled {
			t.Error("failureFn should be called")
		}

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MatchThen should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		err := errors.New("test error")
		var sideEffect string

		result := maybe.Failed[int](err).
			MatchThen(
				func(x int) { sideEffect = "some" },
				func() { sideEffect = "none" },
				func(e error) { sideEffect = "Processing error" },
			).
			Map(func(x int) int { return x * 2 })

		if sideEffect != "Processing error" {
			t.Errorf("expected 'Processing error', got %s", sideEffect)
		}

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("chained operations should return Failure")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("preserves error through multiple MatchThen calls", func(t *testing.T) {
		err := errors.New("persistent error")
		var log []string

		result := maybe.Failed[int](err).
			MatchThen(
				func(x int) { log = append(log, "some") },
				func() { log = append(log, "none") },
				func(e error) { log = append(log, "first") },
			).
			MatchThen(
				func(x int) { log = append(log, "some") },
				func() { log = append(log, "none") },
				func(e error) { log = append(log, "second") },
			)

		if len(log) != 2 || log[0] != "first" || log[1] != "second" {
			t.Errorf("expected [first second], got %v", log)
		}

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("chained MatchThen should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("preserves error state", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		result := failure.MatchThen(
			func(x int) {},
			func() {},
			func(e error) { /* no-op */ },
		)

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MatchThen should return Failure type")
		}
		_, _, gotErr := resultFailure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})
}


func TestFailure_MapIfEmpty(t *testing.T) {
	t.Run("returns original Failure unchanged", func(t *testing.T) {
		originalErr := errors.New("original error")
		failure := maybe.Failed[int](originalErr)
		called := false

		result := failure.MapIfEmpty(func() (int, error) {
			called = true
			return 100, nil
		})

		if called {
			t.Error("MapIfEmpty should not call function for Failure")
		}

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MapIfEmpty should return Failure for Failure")
		}
		_, _, err := resultFailure.Get()
		if err != originalErr {
			t.Errorf("expected %v, got %v", originalErr, err)
		}
	})
}

func TestFailure_MapIfFailed(t *testing.T) {
	t.Run("executes recovery function and returns Some", func(t *testing.T) {
		originalErr := errors.New("not found")
		failure := maybe.Failed[int](originalErr)

		result := failure.MapIfFailed(func(err error) (int, error) {
			if err.Error() == "not found" {
				return 42, nil
			}
			return 0, err
		})

		some, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("MapIfFailed should return Some when recovery succeeds")
		}
		value, _, _ := some.Get()
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})

	t.Run("returns new Failure when recovery returns error", func(t *testing.T) {
		originalErr := errors.New("original error")
		newErr := errors.New("new error")
		failure := maybe.Failed[int](originalErr)

		result := failure.MapIfFailed(func(err error) (int, error) {
			return 0, newErr
		})

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MapIfFailed should return Failure when recovery returns error")
		}
		_, _, err := resultFailure.Get()
		if err != newErr {
			t.Errorf("expected %v, got %v", newErr, err)
		}
	})

	t.Run("transforms error by wrapping with context", func(t *testing.T) {
		dbErr := errors.New("connection timeout")
		failure := maybe.Failed[int](dbErr)

		result := failure.MapIfFailed(func(err error) (int, error) {
			return 0, errors.New("database error: " + err.Error())
		})

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MapIfFailed should return Failure for error transformation")
		}
		_, _, err := resultFailure.Get()
		if err == nil || err.Error() != "database error: connection timeout" {
			t.Errorf("expected wrapped error, got %v", err)
		}
	})

	t.Run("receives original error in function", func(t *testing.T) {
		originalErr := errors.New("test error")
		failure := maybe.Failed[int](originalErr)
		var receivedErr error

		failure.MapIfFailed(func(err error) (int, error) {
			receivedErr = err
			return 0, nil
		})

		if receivedErr != originalErr {
			t.Errorf("expected to receive %v, got %v", originalErr, receivedErr)
		}
	})

	t.Run("catches panic in recovery function", func(t *testing.T) {
		failure := maybe.Failed[int](errors.New("error"))

		result := failure.MapIfFailed(func(err error) (int, error) {
			panic("recovery panic")
		})

		resultFailure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MapIfFailed should return Failure when recovery panics")
		}
		_, _, err := resultFailure.Get()
		if err == nil {
			t.Error("expected error from panic")
		}
	})
}
