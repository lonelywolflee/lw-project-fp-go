package maybe_test

import (
	"errors"
	"testing"

	"github.com/lonelywolflee/lw-project-fp-go/maybe"
)

func TestNone_Get(t *testing.T) {
	t.Run("returns zero value and nil error for int", func(t *testing.T) {
		none := maybe.Empty[int]()
		value, err := none.Get()
		if err != nil {
			t.Errorf("None.Get() should return nil error, got %v", err)
		}
		if value != 0 {
			t.Errorf("expected zero value (0), got %d", value)
		}
	})

	t.Run("returns zero value and nil error for string", func(t *testing.T) {
		none := maybe.Empty[string]()
		value, err := none.Get()
		if err != nil {
			t.Error("None.Get() should return nil error")
		}
		if value != "" {
			t.Errorf("expected zero value (empty string), got %s", value)
		}
	})

	t.Run("returns zero value for bool", func(t *testing.T) {
		none := maybe.Empty[bool]()
		value, err := none.Get()
		if err != nil {
			t.Errorf("None.Get() should return nil error, got %v", err)
		}
		if value != false {
			t.Errorf("expected zero value (false), got %v", value)
		}
	})

	t.Run("returns nil and nil error for pointer type", func(t *testing.T) {
		none := maybe.Empty[*int]()
		value, err := none.Get()
		if err != nil {
			t.Errorf("None.Get() should return nil error, got %v", err)
		}
		if value != nil {
			t.Errorf("expected nil pointer, got %v", value)
		}
	})

	t.Run("returns zero value for struct", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		none := maybe.Empty[User]()
		value, err := none.Get()
		if err != nil {
			t.Errorf("None.Get() should return nil error, got %v", err)
		}
		if value.Name != "" || value.Age != 0 {
			t.Errorf("expected zero value User{}, got %+v", value)
		}
	})
}

func TestNone_Map(t *testing.T) {
	t.Run("returns Empty and ignores function", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.Map(func(x int) int { return x * 2 })

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("None.Map should return None type")
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		none := maybe.Empty[int]()
		executed := false
		result := none.Map(func(x int) int {
			executed = true
			return x * 2
		})

		if executed {
			t.Error("None.Map should not execute the function")
		}

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("None.Map should return None type")
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.Map(func(x int) int {
			panic("this should never be called")
		})

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("None.Map should return None type without executing function")
		}
	})
}

func TestNone_FlatMap(t *testing.T) {
	t.Run("returns Empty and ignores function", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.FlatMap(func(x int) maybe.Maybe[int] {
			return maybe.Just(x * 2)
		})

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("None.FlatMap should return None type")
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		none := maybe.Empty[int]()
		executed := false
		result := none.FlatMap(func(x int) maybe.Maybe[int] {
			executed = true
			return maybe.Just(x * 2)
		})

		if executed {
			t.Error("None.FlatMap should not execute the function")
		}

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("None.FlatMap should return None type")
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		none := maybe.Empty[string]()
		result := none.FlatMap(func(x string) maybe.Maybe[string] {
			panic("this should never be called")
		})

		_, ok := result.(maybe.None[string])
		if !ok {
			t.Fatal("None.FlatMap should return None type without executing function")
		}
	})

	t.Run("chains with Map operation", func(t *testing.T) {
		result := maybe.Empty[int]().
			Map(func(x int) int { return x * 2 })

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("chained operations on None should return None type")
		}
	})
}

func TestNone_Filter(t *testing.T) {
	t.Run("returns None and ignores predicate", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.Filter(func(x int) bool { return x > 5 })

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("None.Filter should return None type")
		}
	})

	t.Run("does not execute the predicate function", func(t *testing.T) {
		none := maybe.Empty[int]()
		executed := false
		result := none.Filter(func(x int) bool {
			executed = true
			return true
		})

		if executed {
			t.Error("None.Filter should not execute the predicate function")
		}

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("None.Filter should return None type")
		}
	})

	t.Run("predicate can panic but never called", func(t *testing.T) {
		none := maybe.Empty[string]()
		result := none.Filter(func(x string) bool {
			panic("this should never be called")
		})

		_, ok := result.(maybe.None[string])
		if !ok {
			t.Fatal("None.Filter should return None type without executing predicate")
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		result := maybe.Empty[int]().
			Filter(func(x int) bool { return x > 5 }).
			Map(func(x int) int { return x * 2 })

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("chained Filter and Map on None should return None type")
		}
	})

	t.Run("preserves None through Filter", func(t *testing.T) {
		result := maybe.Empty[int]().
			Filter(func(x int) bool { return x > 10 })

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("Filter on None should preserve None type")
		}
	})
}

func TestNone_Then(t *testing.T) {
	t.Run("returns None and ignores function", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.Then(func(x int) { /* no-op */ })

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("None.Then should return None type")
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		none := maybe.Empty[int]()
		executed := false
		result := none.Then(func(x int) {
			executed = true
		})

		if executed {
			t.Error("None.Then should not execute the function")
		}

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("None.Then should return None type")
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		none := maybe.Empty[string]()
		result := none.Then(func(x string) {
			panic("this should never be called")
		})

		_, ok := result.(maybe.None[string])
		if !ok {
			t.Fatal("None.Then should return None type without executing function")
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		var sideEffect int

		result := maybe.Empty[int]().
			Then(func(x int) { sideEffect = x }).
			Map(func(x int) int { return x * 2 })

		if sideEffect != 0 {
			t.Errorf("side effect should not be triggered, got %d", sideEffect)
		}

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("chained Then and Map on None should return None type")
		}
	})

	t.Run("preserves None through multiple Then calls", func(t *testing.T) {
		var callCount int

		result := maybe.Empty[int]().
			Then(func(x int) { callCount++ }).
			Then(func(x int) { callCount++ }).
			Then(func(x int) { callCount++ })

		if callCount != 0 {
			t.Errorf("no Then calls should execute, got %d calls", callCount)
		}

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("multiple Then calls on None should preserve None type")
		}
	})
}

func TestNone_OrElseGet(t *testing.T) {
	t.Run("calls function and returns result", func(t *testing.T) {
		none := maybe.Empty[int]()
		called := false
		var receivedErr error
		result := none.OrElseGet(func(err error) int {
			called = true
			receivedErr = err
			return 42
		})

		if !called {
			t.Error("OrElseGet should call the function when None has no value")
		}
		if receivedErr != nil {
			t.Errorf("None should pass nil error, got %v", receivedErr)
		}
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("returns string from function", func(t *testing.T) {
		none := maybe.Empty[string]()
		result := none.OrElseGet(func(err error) string { return "default" })

		if result != "default" {
			t.Errorf("expected 'default', got %s", result)
		}
	})

	t.Run("executes function every time", func(t *testing.T) {
		none := maybe.Empty[int]()
		callCount := 0

		result1 := none.OrElseGet(func(err error) int {
			callCount++
			return callCount
		})
		result2 := none.OrElseGet(func(err error) int {
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
		none := maybe.Empty[[]int]()
		result := none.OrElseGet(func(err error) []int { return []int{1, 2, 3} })

		if len(result) != 3 {
			t.Errorf("expected slice length 3, got %d", len(result))
		}
		if result[0] != 1 || result[1] != 2 || result[2] != 3 {
			t.Errorf("expected [1 2 3], got %v", result)
		}
	})

	t.Run("can return zero values", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.OrElseGet(func(err error) int { return 0 })

		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("function can compute complex values", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.OrElseGet(func(err error) int {
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

	t.Run("receives nil error parameter", func(t *testing.T) {
		none := maybe.Empty[int]()
		var capturedErr error
		result := none.OrElseGet(func(err error) int {
			capturedErr = err
			return 100
		})

		if capturedErr != nil {
			t.Errorf("expected nil error for None, got %v", capturedErr)
		}
		if result != 100 {
			t.Errorf("expected 100, got %d", result)
		}
	})
}

func TestNone_OrElseDefault(t *testing.T) {
	t.Run("returns the default value", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.OrElseDefault(42)

		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("returns string default value", func(t *testing.T) {
		none := maybe.Empty[string]()
		result := none.OrElseDefault("default")

		if result != "default" {
			t.Errorf("expected 'default', got %s", result)
		}
	})

	t.Run("works with different types", func(t *testing.T) {
		none := maybe.Empty[[]int]()
		result := none.OrElseDefault([]int{1, 2, 3})

		if len(result) != 3 {
			t.Errorf("expected slice length 3, got %d", len(result))
		}
		if result[0] != 1 || result[1] != 2 || result[2] != 3 {
			t.Errorf("expected [1 2 3], got %v", result)
		}
	})

	t.Run("can return zero values", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.OrElseDefault(0)

		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("returns same default every time", func(t *testing.T) {
		none := maybe.Empty[int]()
		result1 := none.OrElseDefault(10)
		result2 := none.OrElseDefault(10)

		if result1 != 10 || result2 != 10 {
			t.Errorf("expected both results to be 10, got %d and %d", result1, result2)
		}
	})
}

func TestNone_MatchThen(t *testing.T) {
	t.Run("executes noneFn and returns original None", func(t *testing.T) {
		none := maybe.Empty[int]()
		someCalled := false
		noneCalled := false
		failureCalled := false

		result := none.MatchThen(
			func(x int) { someCalled = true },
			func() { noneCalled = true },
			func(err error) { failureCalled = true },
		)

		if someCalled {
			t.Error("someFn should not be called")
		}
		if !noneCalled {
			t.Error("noneFn should be called")
		}
		if failureCalled {
			t.Error("failureFn should not be called")
		}

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("MatchThen should return None type")
		}
	})

	t.Run("can be used for logging", func(t *testing.T) {
		none := maybe.Empty[string]()
		var log string

		result := none.MatchThen(
			func(x string) { log = "Got value: " + x },
			func() { log = "No value" },
			func(err error) { log = "Error: " + err.Error() },
		)

		if log != "No value" {
			t.Errorf("expected 'No value', got %s", log)
		}

		_, ok := result.(maybe.None[string])
		if !ok {
			t.Fatal("MatchThen should return None type")
		}
	})

	t.Run("catches panic in noneFn and converts to Failure", func(t *testing.T) {
		none := maybe.Empty[int]()

		result := none.MatchThen(
			func(x int) {},
			func() { panic("noneFn panic") },
			func(err error) {},
		)

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MatchThen should return Failure when noneFn panics")
		}
		_, gotErr := failure.Get()
		if gotErr.Error() != "noneFn panic" {
			t.Errorf("expected panic message, got %s", gotErr.Error())
		}
	})

	t.Run("catches panic with error type in noneFn", func(t *testing.T) {
		none := maybe.Empty[int]()
		testErr := errors.New("test error")

		result := none.MatchThen(
			func(x int) {},
			func() { panic(testErr) },
			func(err error) {},
		)

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MatchThen should return Failure when noneFn panics")
		}
		_, gotErr := failure.Get()
		if gotErr != testErr {
			t.Errorf("expected %v, got %v", testErr, gotErr)
		}
	})

	t.Run("someFn and failureFn can panic but never called", func(t *testing.T) {
		none := maybe.Empty[int]()
		noneCalled := false

		result := none.MatchThen(
			func(x int) { panic("someFn should not be called") },
			func() { noneCalled = true },
			func(err error) { panic("failureFn should not be called") },
		)

		if !noneCalled {
			t.Error("noneFn should be called")
		}

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("MatchThen should return None type")
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		var sideEffect string

		result := maybe.Empty[int]().
			MatchThen(
				func(x int) { sideEffect = "some" },
				func() { sideEffect = "Processing None" },
				func(err error) { sideEffect = "error" },
			).
			Map(func(x int) int { return x * 2 })

		if sideEffect != "Processing None" {
			t.Errorf("expected 'Processing None', got %s", sideEffect)
		}

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("chained operations should return None")
		}
	})

	t.Run("can be chained multiple times", func(t *testing.T) {
		var log []string

		result := maybe.Empty[int]().
			MatchThen(
				func(x int) { log = append(log, "some") },
				func() { log = append(log, "first") },
				func(err error) { log = append(log, "error") },
			).
			MatchThen(
				func(x int) { log = append(log, "some") },
				func() { log = append(log, "second") },
				func(err error) { log = append(log, "error") },
			)

		if len(log) != 2 || log[0] != "first" || log[1] != "second" {
			t.Errorf("expected [first second], got %v", log)
		}

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("chained MatchThen should return None type")
		}
	})

	t.Run("preserves None state", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.MatchThen(
			func(x int) {},
			func() { /* no-op */ },
			func(err error) {},
		)

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("MatchThen should return None type")
		}
	})
}


func TestNone_MapIfEmpty(t *testing.T) {
	t.Run("executes recovery function and returns Some", func(t *testing.T) {
		none := maybe.Empty[int]()

		result := none.MapIfEmpty(func() (int, error) {
			return 42, nil
		})

		some, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("MapIfEmpty should return Some when recovery succeeds")
		}
		value, _ := some.Get()
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})

	t.Run("returns Failure when recovery function returns error", func(t *testing.T) {
		none := maybe.Empty[int]()
		testErr := errors.New("recovery failed")

		result := none.MapIfEmpty(func() (int, error) {
			return 0, testErr
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MapIfEmpty should return Failure when recovery returns error")
		}
		_, err := failure.Get()
		if err != testErr {
			t.Errorf("expected error %v, got %v", testErr, err)
		}
	})

	t.Run("transforms None to Failure with custom error", func(t *testing.T) {
		none := maybe.Empty[int]()
		customErr := errors.New("value required")

		result := none.MapIfEmpty(func() (int, error) {
			return 0, customErr
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MapIfEmpty should return Failure for error transformation")
		}
		_, err := failure.Get()
		if err != customErr {
			t.Errorf("expected error %v, got %v", customErr, err)
		}
	})

	t.Run("catches panic in recovery function", func(t *testing.T) {
		none := maybe.Empty[int]()

		result := none.MapIfEmpty(func() (int, error) {
			panic("something went wrong")
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MapIfEmpty should return Failure when recovery panics")
		}
		_, err := failure.Get()
		if err == nil {
			t.Error("expected error from panic")
		}
	})
}

func TestNone_MapIfFailed(t *testing.T) {
	t.Run("returns original None unchanged", func(t *testing.T) {
		none := maybe.Empty[int]()
		called := false

		result := none.MapIfFailed(func(err error) (int, error) {
			called = true
			return 100, nil
		})

		if called {
			t.Error("MapIfFailed should not call function for None")
		}

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("MapIfFailed should return None for None")
		}
	})
}
