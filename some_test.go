package lwfp

import (
	"errors"
	"testing"
)

func TestSome_GetValue(t *testing.T) {
	t.Run("returns the wrapped value", func(t *testing.T) {
		some := Some[int]{v: 42}
		if some.GetValue() != 42 {
			t.Errorf("expected 42, got %d", some.GetValue())
		}
	})

	t.Run("returns string value", func(t *testing.T) {
		some := Some[string]{v: "test"}
		if some.GetValue() != "test" {
			t.Errorf("expected 'test', got %s", some.GetValue())
		}
	})
}

func TestSome_Map(t *testing.T) {
	t.Run("transforms value successfully", func(t *testing.T) {
		some := Just(5)
		result := some.Map(func(x int) any { return x * 2 })

		resultSome, ok := result.(Some[any])
		if !ok {
			t.Fatal("Map should return Some type")
		}
		if resultSome.GetValue() != 10 {
			t.Errorf("expected 10, got %v", resultSome.GetValue())
		}
	})

	t.Run("handles string transformation", func(t *testing.T) {
		some := Just("hello")
		result := some.Map(func(x string) any { return x + " world" })

		resultSome, ok := result.(Some[any])
		if !ok {
			t.Fatal("Map should return Some type")
		}
		if resultSome.GetValue() != "hello world" {
			t.Errorf("expected 'hello world', got %v", resultSome.GetValue())
		}
	})

	t.Run("catches panic and converts to Failure", func(t *testing.T) {
		some := Just(5)
		result := some.Map(func(x int) any {
			panic("something went wrong")
		})

		failure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("Map should return Failure when panic occurs")
		}
		if failure.GetError().Error() != "something went wrong" {
			t.Errorf("expected panic message, got %s", failure.GetError().Error())
		}
	})

	t.Run("catches panic with error type", func(t *testing.T) {
		some := Just(5)
		testErr := errors.New("test error")
		result := some.Map(func(x int) any {
			panic(testErr)
		})

		failure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("Map should return Failure when panic occurs")
		}
		if failure.GetError() != testErr {
			t.Errorf("expected %v, got %v", testErr, failure.GetError())
		}
	})

	t.Run("handles nil pointer dereference panic", func(t *testing.T) {
		some := Just(10)
		result := some.Map(func(x int) any {
			var ptr *int
			return *ptr // This will panic
		})

		_, ok := result.(Failure[any])
		if !ok {
			t.Fatal("Map should return Failure when panic occurs")
		}
	})
}

func TestSome_FlatMap(t *testing.T) {
	t.Run("chains Maybe values successfully", func(t *testing.T) {
		some := Just(5)
		result := some.FlatMap(func(x int) Maybe[any] {
			return Just[any](x * 2)
		})

		resultSome, ok := result.(Some[any])
		if !ok {
			t.Fatal("FlatMap should return Some type")
		}
		if resultSome.GetValue() != 10 {
			t.Errorf("expected 10, got %v", resultSome.GetValue())
		}
	})

	t.Run("returns Empty when function returns Empty", func(t *testing.T) {
		some := Just(5)
		result := some.FlatMap(func(x int) Maybe[any] {
			if x < 10 {
				return Empty[any]()
			}
			return Just[any](x)
		})

		_, ok := result.(None[any])
		if !ok {
			t.Fatal("FlatMap should return None when function returns Empty")
		}
	})

	t.Run("returns Failure when function returns Failure", func(t *testing.T) {
		some := Just(5)
		testErr := errors.New("validation error")
		result := some.FlatMap(func(x int) Maybe[any] {
			if x < 10 {
				return Fail[any](testErr)
			}
			return Just[any](x)
		})

		failure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("FlatMap should return Failure when function returns Failure")
		}
		if failure.GetError() != testErr {
			t.Errorf("expected %v, got %v", testErr, failure.GetError())
		}
	})

	t.Run("catches panic and converts to Failure", func(t *testing.T) {
		some := Just(5)
		result := some.FlatMap(func(x int) Maybe[any] {
			panic("flatmap panic")
		})

		failure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("FlatMap should return Failure when panic occurs")
		}
		if failure.GetError().Error() != "flatmap panic" {
			t.Errorf("expected panic message, got %s", failure.GetError().Error())
		}
	})

	t.Run("catches panic with error type", func(t *testing.T) {
		some := Just(5)
		testErr := errors.New("flatmap error")
		result := some.FlatMap(func(x int) Maybe[any] {
			panic(testErr)
		})

		failure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("FlatMap should return Failure when panic occurs")
		}
		if failure.GetError() != testErr {
			t.Errorf("expected %v, got %v", testErr, failure.GetError())
		}
	})

	t.Run("chains FlatMap with conditional logic", func(t *testing.T) {
		result := Just(15).FlatMap(func(x int) Maybe[any] {
			if x > 10 {
				return Just[any](x * 2)
			}
			return Empty[any]()
		})

		resultSome, ok := result.(Some[any])
		if !ok {
			t.Fatal("FlatMap should return Some type")
		}
		if resultSome.GetValue() != 30 {
			t.Errorf("expected 30, got %v", resultSome.GetValue())
		}
	})
}

func TestSome_Filter(t *testing.T) {
	t.Run("returns Some when predicate is true", func(t *testing.T) {
		some := Just(10)
		result := some.Filter(func(x int) bool { return x > 5 })

		resultSome, ok := result.(Some[int])
		if !ok {
			t.Fatal("Filter should return Some when predicate is true")
		}
		if resultSome.GetValue() != 10 {
			t.Errorf("expected 10, got %d", resultSome.GetValue())
		}
	})

	t.Run("returns None when predicate is false", func(t *testing.T) {
		some := Just(3)
		result := some.Filter(func(x int) bool { return x > 5 })

		_, ok := result.(None[int])
		if !ok {
			t.Fatal("Filter should return None when predicate is false")
		}
	})

	t.Run("handles string filtering", func(t *testing.T) {
		some := Just("hello")
		result := some.Filter(func(x string) bool { return len(x) > 3 })

		resultSome, ok := result.(Some[string])
		if !ok {
			t.Fatal("Filter should return Some when predicate is true")
		}
		if resultSome.GetValue() != "hello" {
			t.Errorf("expected 'hello', got %s", resultSome.GetValue())
		}
	})

	t.Run("filters out short strings", func(t *testing.T) {
		some := Just("hi")
		result := some.Filter(func(x string) bool { return len(x) > 3 })

		_, ok := result.(None[string])
		if !ok {
			t.Fatal("Filter should return None when predicate is false")
		}
	})

	t.Run("catches panic and converts to Failure", func(t *testing.T) {
		some := Just(5)
		result := some.Filter(func(x int) bool {
			panic("filter panic")
		})

		failure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Filter should return Failure when panic occurs")
		}
		if failure.GetError().Error() != "filter panic" {
			t.Errorf("expected panic message, got %s", failure.GetError().Error())
		}
	})

	t.Run("catches panic with error type", func(t *testing.T) {
		some := Just(5)
		testErr := errors.New("filter error")
		result := some.Filter(func(x int) bool {
			panic(testErr)
		})

		failure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Filter should return Failure when panic occurs")
		}
		if failure.GetError() != testErr {
			t.Errorf("expected %v, got %v", testErr, failure.GetError())
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		result := Just(10).
			Filter(func(x int) bool { return x > 5 }).
			Map(func(x int) any { return x * 2 })

		resultSome, ok := result.(Some[any])
		if !ok {
			t.Fatal("chained Filter and Map should return Some")
		}
		if resultSome.GetValue() != 20 {
			t.Errorf("expected 20, got %v", resultSome.GetValue())
		}
	})

	t.Run("returns None when filter fails in chain", func(t *testing.T) {
		result := Just(3).
			Filter(func(x int) bool { return x > 5 }).
			Map(func(x int) any { return x * 2 })

		_, ok := result.(None[any])
		if !ok {
			t.Fatal("chained Filter should return None when predicate is false")
		}
	})

	t.Run("handles complex predicate", func(t *testing.T) {
		some := Just(42)
		result := some.Filter(func(x int) bool {
			return x%2 == 0 && x > 10
		})

		resultSome, ok := result.(Some[int])
		if !ok {
			t.Fatal("Filter should return Some when complex predicate is true")
		}
		if resultSome.GetValue() != 42 {
			t.Errorf("expected 42, got %d", resultSome.GetValue())
		}
	})
}

func TestSome_Then(t *testing.T) {
	t.Run("executes function and returns original Some", func(t *testing.T) {
		executed := false
		var capturedValue int

		some := Just(10)
		result := some.Then(func(x int) {
			executed = true
			capturedValue = x
		})

		if !executed {
			t.Error("Then should execute the function")
		}
		if capturedValue != 10 {
			t.Errorf("expected captured value 10, got %d", capturedValue)
		}

		resultSome, ok := result.(Some[int])
		if !ok {
			t.Fatal("Then should return Some type")
		}
		if resultSome.GetValue() != 10 {
			t.Errorf("expected 10, got %d", resultSome.GetValue())
		}
	})

	t.Run("can be used for logging", func(t *testing.T) {
		var loggedValue string

		result := Just("hello").Then(func(x string) {
			loggedValue = "Logged: " + x
		})

		if loggedValue != "Logged: hello" {
			t.Errorf("expected 'Logged: hello', got %s", loggedValue)
		}

		resultSome, ok := result.(Some[string])
		if !ok {
			t.Fatal("Then should return Some type")
		}
		if resultSome.GetValue() != "hello" {
			t.Errorf("expected 'hello', got %s", resultSome.GetValue())
		}
	})

	t.Run("catches panic and converts to Failure", func(t *testing.T) {
		some := Just(5)
		result := some.Then(func(x int) {
			panic("then panic")
		})

		failure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Then should return Failure when panic occurs")
		}
		if failure.GetError().Error() != "then panic" {
			t.Errorf("expected panic message, got %s", failure.GetError().Error())
		}
	})

	t.Run("catches panic with error type", func(t *testing.T) {
		some := Just(5)
		testErr := errors.New("then error")
		result := some.Then(func(x int) {
			panic(testErr)
		})

		failure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Then should return Failure when panic occurs")
		}
		if failure.GetError() != testErr {
			t.Errorf("expected %v, got %v", testErr, failure.GetError())
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		var sideEffect int

		result := Just(10).
			Then(func(x int) { sideEffect = x }).
			Map(func(x int) any { return x * 2 })

		if sideEffect != 10 {
			t.Errorf("expected side effect to be 10, got %d", sideEffect)
		}

		resultSome, ok := result.(Some[any])
		if !ok {
			t.Fatal("chained Then and Map should return Some")
		}
		if resultSome.GetValue() != 20 {
			t.Errorf("expected 20, got %v", resultSome.GetValue())
		}
	})

	t.Run("can be chained multiple times", func(t *testing.T) {
		var log []string

		result := Just(5).
			Then(func(x int) { log = append(log, "first") }).
			Then(func(x int) { log = append(log, "second") }).
			Then(func(x int) { log = append(log, "third") })

		if len(log) != 3 {
			t.Errorf("expected 3 log entries, got %d", len(log))
		}
		if log[0] != "first" || log[1] != "second" || log[2] != "third" {
			t.Errorf("expected [first second third], got %v", log)
		}

		resultSome, ok := result.(Some[int])
		if !ok {
			t.Fatal("chained Then should return Some type")
		}
		if resultSome.GetValue() != 5 {
			t.Errorf("expected 5, got %d", resultSome.GetValue())
		}
	})

	t.Run("does not change value", func(t *testing.T) {
		original := Just(42)
		result := original.Then(func(x int) {
			// Try to modify (won't affect original)
			x = 100
		})

		resultSome, ok := result.(Some[int])
		if !ok {
			t.Fatal("Then should return Some type")
		}
		if resultSome.GetValue() != 42 {
			t.Errorf("value should remain 42, got %d", resultSome.GetValue())
		}
	})
}

func TestSome_OrElseGet(t *testing.T) {
	t.Run("returns the value and does not call function", func(t *testing.T) {
		some := Just(42)
		called := false
		result := some.OrElseGet(func() int {
			called = true
			return 0
		})

		if called {
			t.Error("OrElseGet should not call the function when Some has a value")
		}
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("returns string value without calling function", func(t *testing.T) {
		some := Just("hello")
		result := some.OrElseGet(func() string { return "default" })

		if result != "hello" {
			t.Errorf("expected 'hello', got %s", result)
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		some := Just(10)
		result := some.OrElseGet(func() int {
			panic("this should never be called")
		})

		if result != 10 {
			t.Errorf("expected 10, got %d", result)
		}
	})

	t.Run("works with different types", func(t *testing.T) {
		some := Just([]int{1, 2, 3})
		result := some.OrElseGet(func() []int { return []int{} })

		if len(result) != 3 {
			t.Errorf("expected slice length 3, got %d", len(result))
		}
		if result[0] != 1 || result[1] != 2 || result[2] != 3 {
			t.Errorf("expected [1 2 3], got %v", result)
		}
	})

	t.Run("works with zero values", func(t *testing.T) {
		some := Just(0)
		result := some.OrElseGet(func() int { return 42 })

		if result != 0 {
			t.Errorf("expected 0 (the actual value), got %d", result)
		}
	})
}

func TestSome_OrElseDefault(t *testing.T) {
	t.Run("returns the value and ignores default", func(t *testing.T) {
		some := Just(42)
		result := some.OrElseDefault(0)

		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("returns string value ignoring default", func(t *testing.T) {
		some := Just("hello")
		result := some.OrElseDefault("default")

		if result != "hello" {
			t.Errorf("expected 'hello', got %s", result)
		}
	})

	t.Run("works with different types", func(t *testing.T) {
		some := Just([]int{1, 2, 3})
		result := some.OrElseDefault([]int{})

		if len(result) != 3 {
			t.Errorf("expected slice length 3, got %d", len(result))
		}
		if result[0] != 1 || result[1] != 2 || result[2] != 3 {
			t.Errorf("expected [1 2 3], got %v", result)
		}
	})

	t.Run("works with zero values", func(t *testing.T) {
		some := Just(0)
		result := some.OrElseDefault(42)

		if result != 0 {
			t.Errorf("expected 0 (the actual value), got %d", result)
		}
	})

	t.Run("can use same value as default", func(t *testing.T) {
		some := Just(10)
		result := some.OrElseDefault(10)

		if result != 10 {
			t.Errorf("expected 10, got %d", result)
		}
	})
}
