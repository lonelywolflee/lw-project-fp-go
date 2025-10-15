package maybe_test

import (
	"errors"
	"testing"

	"github.com/lonelywolflee/lw-project-fp-go/maybe"
)

func TestNone_Get(t *testing.T) {
	t.Run("returns nil", func(t *testing.T) {
		none := maybe.Empty[int]()
		if none.Get() != nil {
			t.Error("None.Get() should return nil")
		}
	})

	t.Run("returns nil for different types", func(t *testing.T) {
		none := maybe.Empty[string]()
		if none.Get() != nil {
			t.Error("None.Get() should return nil")
		}
	})
}

func TestNone_Map(t *testing.T) {
	t.Run("returns Empty and ignores function", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.Map(func(x int) any { return x * 2 })

		_, ok := result.(maybe.None[any])
		if !ok {
			t.Fatal("None.Map should return None type")
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		none := maybe.Empty[int]()
		executed := false
		result := none.Map(func(x int) any {
			executed = true
			return x * 2
		})

		if executed {
			t.Error("None.Map should not execute the function")
		}

		_, ok := result.(maybe.None[any])
		if !ok {
			t.Fatal("None.Map should return None type")
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.Map(func(x int) any {
			panic("this should never be called")
		})

		_, ok := result.(maybe.None[any])
		if !ok {
			t.Fatal("None.Map should return None type without executing function")
		}
	})
}

func TestNone_FlatMap(t *testing.T) {
	t.Run("returns Empty and ignores function", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.FlatMap(func(x int) maybe.Maybe[any] {
			return maybe.Just[any](x * 2)
		})

		_, ok := result.(maybe.None[any])
		if !ok {
			t.Fatal("None.FlatMap should return None type")
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		none := maybe.Empty[int]()
		executed := false
		result := none.FlatMap(func(x int) maybe.Maybe[any] {
			executed = true
			return maybe.Just[any](x * 2)
		})

		if executed {
			t.Error("None.FlatMap should not execute the function")
		}

		_, ok := result.(maybe.None[any])
		if !ok {
			t.Fatal("None.FlatMap should return None type")
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		none := maybe.Empty[string]()
		result := none.FlatMap(func(x string) maybe.Maybe[any] {
			panic("this should never be called")
		})

		_, ok := result.(maybe.None[any])
		if !ok {
			t.Fatal("None.FlatMap should return None type without executing function")
		}
	})

	t.Run("chains with Map operation", func(t *testing.T) {
		result := maybe.Empty[int]().
			Map(func(x int) any { return x * 2 })

		_, ok := result.(maybe.None[any])
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
			Map(func(x int) any { return x * 2 })

		_, ok := result.(maybe.None[any])
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
			Map(func(x int) any { return x * 2 })

		if sideEffect != 0 {
			t.Errorf("side effect should not be triggered, got %d", sideEffect)
		}

		_, ok := result.(maybe.None[any])
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
		result := none.OrElseGet(func() int {
			called = true
			return 42
		})

		if !called {
			t.Error("OrElseGet should call the function when None has no value")
		}
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("returns string from function", func(t *testing.T) {
		none := maybe.Empty[string]()
		result := none.OrElseGet(func() string { return "default" })

		if result != "default" {
			t.Errorf("expected 'default', got %s", result)
		}
	})

	t.Run("executes function every time", func(t *testing.T) {
		none := maybe.Empty[int]()
		callCount := 0

		result1 := none.OrElseGet(func() int {
			callCount++
			return callCount
		})
		result2 := none.OrElseGet(func() int {
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
		result := none.OrElseGet(func() []int { return []int{1, 2, 3} })

		if len(result) != 3 {
			t.Errorf("expected slice length 3, got %d", len(result))
		}
		if result[0] != 1 || result[1] != 2 || result[2] != 3 {
			t.Errorf("expected [1 2 3], got %v", result)
		}
	})

	t.Run("can return zero values", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.OrElseGet(func() int { return 0 })

		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("function can compute complex values", func(t *testing.T) {
		none := maybe.Empty[int]()
		result := none.OrElseGet(func() int {
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
		if failure.GetError().Error() != "noneFn panic" {
			t.Errorf("expected panic message, got %s", failure.GetError().Error())
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
		if failure.GetError() != testErr {
			t.Errorf("expected %v, got %v", testErr, failure.GetError())
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
			Map(func(x int) any { return x * 2 })

		if sideEffect != "Processing None" {
			t.Errorf("expected 'Processing None', got %s", sideEffect)
		}

		_, ok := result.(maybe.None[any])
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
