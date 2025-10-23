package maybe_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lonelywolflee/lw-project-fp-go/maybe"
)

// TestSome_GetValue is removed - use TestSome_Get instead
// Get() is now the unified interface method for accessing values

func TestSome_Get(t *testing.T) {
	t.Run("returns value and nil error", func(t *testing.T) {
		some := maybe.Just(42)
		value, err := some.Get()

		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})

	t.Run("returns string value and nil error", func(t *testing.T) {
		some := maybe.Just("hello")
		value, err := some.Get()

		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if value != "hello" {
			t.Errorf("expected 'hello', got %s", value)
		}
	})

	t.Run("works with zero values", func(t *testing.T) {
		some := maybe.Just(0)
		value, err := some.Get()

		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if value != 0 {
			t.Errorf("expected 0, got %d", value)
		}
	})

	t.Run("works with complex types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		user := User{Name: "Alice", Age: 30}
		some := maybe.Just(user)
		value, err := some.Get()

		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if value.Name != "Alice" || value.Age != 30 {
			t.Errorf("expected User{Alice, 30}, got %+v", value)
		}
	})

	t.Run("works with pointers", func(t *testing.T) {
		num := 42
		some := maybe.Just(&num)
		value, err := some.Get()

		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if value == nil || *value != 42 {
			t.Errorf("expected pointer to 42, got %v", value)
		}
	})
}

func TestSome_Map(t *testing.T) {
	t.Run("transforms value successfully", func(t *testing.T) {
		some := maybe.Just(5)
		result := some.Map(func(x int) int { return x * 2 })

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("Map should return Some type")
		}
		value, _ := resultSome.Get()
		if value != 10 {
			t.Errorf("expected 10, got %v", value)
		}
	})

	t.Run("handles string transformation", func(t *testing.T) {
		some := maybe.Just("hello")
		result := some.Map(func(x string) string { return x + " world" })

		resultSome, ok := result.(maybe.Some[string])
		if !ok {
			t.Fatal("Map should return Some type")
		}
		value, _ := resultSome.Get()
		if value != "hello world" {
			t.Errorf("expected 'hello world', got %v", value)
		}
	})

	t.Run("catches panic and converts to Failure", func(t *testing.T) {
		some := maybe.Just(5)
		result := some.Map(func(x int) int {
			panic("something went wrong")
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Map should return Failure when panic occurs")
		}
		_, err := failure.Get()
		if err.Error() != "something went wrong" {
			t.Errorf("expected panic message, got %s", err.Error())
		}
	})

	t.Run("catches panic with error type", func(t *testing.T) {
		some := maybe.Just(5)
		testErr := errors.New("test error")
		result := some.Map(func(x int) int {
			panic(testErr)
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Map should return Failure when panic occurs")
		}
		_, err := failure.Get()
		if err != testErr {
			t.Errorf("expected %v, got %v", testErr, err)
		}
	})

	t.Run("handles nil pointer dereference panic", func(t *testing.T) {
		some := maybe.Just(10)
		result := some.Map(func(x int) int {
			var ptr *int
			return *ptr // This will panic
		})

		_, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Map should return Failure when panic occurs")
		}
	})
}

func TestSome_FlatMap(t *testing.T) {
	t.Run("chains Maybe values successfully", func(t *testing.T) {
		some := maybe.Just(5)
		result := some.FlatMap(func(x int) maybe.Maybe[int] {
			return maybe.Just(x * 2)
		})

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("FlatMap should return Some type")
		}
		value, _ := resultSome.Get()
		if value != 10 {
			t.Errorf("expected 10, got %v", value)
		}
	})

	t.Run("returns Empty when function returns Empty", func(t *testing.T) {
		some := maybe.Just(5)
		result := some.FlatMap(func(x int) maybe.Maybe[int] {
			if x < 10 {
				return maybe.Empty[int]()
			}
			return maybe.Just(x)
		})

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("FlatMap should return None when function returns Empty")
		}
	})

	t.Run("returns Failure when function returns Failure", func(t *testing.T) {
		some := maybe.Just(5)
		testErr := errors.New("validation error")
		result := some.FlatMap(func(x int) maybe.Maybe[int] {
			if x < 10 {
				return maybe.Failed[int](testErr)
			}
			return maybe.Just(x)
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("FlatMap should return Failure when function returns Failure")
		}
		_, err := failure.Get()
		if err != testErr {
			t.Errorf("expected %v, got %v", testErr, err)
		}
	})

	t.Run("catches panic and converts to Failure", func(t *testing.T) {
		some := maybe.Just(5)
		result := some.FlatMap(func(x int) maybe.Maybe[int] {
			panic("flatmap panic")
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("FlatMap should return Failure when panic occurs")
		}
		_, err := failure.Get()
		if err.Error() != "flatmap panic" {
			t.Errorf("expected panic message, got %s", err.Error())
		}
	})

	t.Run("catches panic with error type", func(t *testing.T) {
		some := maybe.Just(5)
		testErr := errors.New("flatmap error")
		result := some.FlatMap(func(x int) maybe.Maybe[int] {
			panic(testErr)
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("FlatMap should return Failure when panic occurs")
		}
		_, err := failure.Get()
		if err != testErr {
			t.Errorf("expected %v, got %v", testErr, err)
		}
	})

	t.Run("chains FlatMap with conditional logic", func(t *testing.T) {
		result := maybe.Just(15).FlatMap(func(x int) maybe.Maybe[int] {
			if x > 10 {
				return maybe.Just(x * 2)
			}
			return maybe.Empty[int]()
		})

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("FlatMap should return Some type")
		}
		value, _ := resultSome.Get()
		if value != 30 {
			t.Errorf("expected 30, got %v", value)
		}
	})
}

func TestSome_Filter(t *testing.T) {
	t.Run("returns Some when predicate is true", func(t *testing.T) {
		some := maybe.Just(10)
		result := some.Filter(func(x int) bool { return x > 5 })

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("Filter should return Some when predicate is true")
		}
		value, _ := resultSome.Get()
		if value != 10 {
			t.Errorf("expected 10, got %d", value)
		}
	})

	t.Run("returns None when predicate is false", func(t *testing.T) {
		some := maybe.Just(3)
		result := some.Filter(func(x int) bool { return x > 5 })

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("Filter should return None when predicate is false")
		}
	})

	t.Run("handles string filtering", func(t *testing.T) {
		some := maybe.Just("hello")
		result := some.Filter(func(x string) bool { return len(x) > 3 })

		resultSome, ok := result.(maybe.Some[string])
		if !ok {
			t.Fatal("Filter should return Some when predicate is true")
		}
		value, _ := resultSome.Get()
		if value != "hello" {
			t.Errorf("expected 'hello', got %s", value)
		}
	})

	t.Run("filters out short strings", func(t *testing.T) {
		some := maybe.Just("hi")
		result := some.Filter(func(x string) bool { return len(x) > 3 })

		_, ok := result.(maybe.None[string])
		if !ok {
			t.Fatal("Filter should return None when predicate is false")
		}
	})

	t.Run("catches panic and converts to Failure", func(t *testing.T) {
		some := maybe.Just(5)
		result := some.Filter(func(x int) bool {
			panic("filter panic")
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Filter should return Failure when panic occurs")
		}
		_, err := failure.Get()
		if err.Error() != "filter panic" {
			t.Errorf("expected panic message, got %s", err.Error())
		}
	})

	t.Run("catches panic with error type", func(t *testing.T) {
		some := maybe.Just(5)
		testErr := errors.New("filter error")
		result := some.Filter(func(x int) bool {
			panic(testErr)
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Filter should return Failure when panic occurs")
		}
		_, err := failure.Get()
		if err != testErr {
			t.Errorf("expected %v, got %v", testErr, err)
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		result := maybe.Just(10).
			Filter(func(x int) bool { return x > 5 }).
			Map(func(x int) int { return x * 2 })

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("chained Filter and Map should return Some")
		}
		value, _ := resultSome.Get()
		if value != 20 {
			t.Errorf("expected 20, got %v", value)
		}
	})

	t.Run("returns None when filter fails in chain", func(t *testing.T) {
		result := maybe.Just(3).
			Filter(func(x int) bool { return x > 5 }).
			Map(func(x int) int { return x * 2 })

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("chained Filter should return None when predicate is false")
		}
	})

	t.Run("handles complex predicate", func(t *testing.T) {
		some := maybe.Just(42)
		result := some.Filter(func(x int) bool {
			return x%2 == 0 && x > 10
		})

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("Filter should return Some when complex predicate is true")
		}
		value, _ := resultSome.Get()
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})
}

func TestSome_Then(t *testing.T) {
	t.Run("executes function and returns original Some", func(t *testing.T) {
		executed := false
		var capturedValue int

		some := maybe.Just(10)
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

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("Then should return Some type")
		}
		value, _ := resultSome.Get()
		if value != 10 {
			t.Errorf("expected 10, got %d", value)
		}
	})

	t.Run("can be used for logging", func(t *testing.T) {
		var loggedValue string

		result := maybe.Just("hello").Then(func(x string) {
			loggedValue = "Logged: " + x
		})

		if loggedValue != "Logged: hello" {
			t.Errorf("expected 'Logged: hello', got %s", loggedValue)
		}

		resultSome, ok := result.(maybe.Some[string])
		if !ok {
			t.Fatal("Then should return Some type")
		}
		value, _ := resultSome.Get()
		if value != "hello" {
			t.Errorf("expected 'hello', got %s", value)
		}
	})

	t.Run("catches panic and converts to Failure", func(t *testing.T) {
		some := maybe.Just(5)
		result := some.Then(func(x int) {
			panic("then panic")
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Then should return Failure when panic occurs")
		}
		_, err := failure.Get()
		if err.Error() != "then panic" {
			t.Errorf("expected panic message, got %s", err.Error())
		}
	})

	t.Run("catches panic with error type", func(t *testing.T) {
		some := maybe.Just(5)
		testErr := errors.New("then error")
		result := some.Then(func(x int) {
			panic(testErr)
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Then should return Failure when panic occurs")
		}
		_, err := failure.Get()
		if err != testErr {
			t.Errorf("expected %v, got %v", testErr, err)
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		var sideEffect int

		result := maybe.Just(10).
			Then(func(x int) { sideEffect = x }).
			Map(func(x int) int { return x * 2 })

		if sideEffect != 10 {
			t.Errorf("expected side effect to be 10, got %d", sideEffect)
		}

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("chained Then and Map should return Some")
		}
		value, _ := resultSome.Get()
		if value != 20 {
			t.Errorf("expected 20, got %v", value)
		}
	})

	t.Run("can be chained multiple times", func(t *testing.T) {
		var log []string

		result := maybe.Just(5).
			Then(func(x int) { log = append(log, "first") }).
			Then(func(x int) { log = append(log, "second") }).
			Then(func(x int) { log = append(log, "third") })

		if len(log) != 3 {
			t.Errorf("expected 3 log entries, got %d", len(log))
		}
		if log[0] != "first" || log[1] != "second" || log[2] != "third" {
			t.Errorf("expected [first second third], got %v", log)
		}

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("chained Then should return Some type")
		}
		value, _ := resultSome.Get()
		if value != 5 {
			t.Errorf("expected 5, got %d", value)
		}
	})

	t.Run("does not change value", func(t *testing.T) {
		original := maybe.Just(42)
		result := original.Then(func(x int) {
			// Try to modify (won't affect original)
			x = 100
		})

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("Then should return Some type")
		}
		value, _ := resultSome.Get()
		if value != 42 {
			t.Errorf("value should remain 42, got %d", value)
		}
	})
}

func TestSome_OrElseGet(t *testing.T) {
	t.Run("returns the value and does not call function", func(t *testing.T) {
		some := maybe.Just(42)
		called := false
		result := some.OrElseGet(func(err error) int {
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
		some := maybe.Just("hello")
		result := some.OrElseGet(func(err error) string { return "default" })

		if result != "hello" {
			t.Errorf("expected 'hello', got %s", result)
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		some := maybe.Just(10)
		result := some.OrElseGet(func(err error) int {
			panic("this should never be called")
		})

		if result != 10 {
			t.Errorf("expected 10, got %d", result)
		}
	})

	t.Run("works with different types", func(t *testing.T) {
		some := maybe.Just([]int{1, 2, 3})
		result := some.OrElseGet(func(err error) []int { return []int{} })

		if len(result) != 3 {
			t.Errorf("expected slice length 3, got %d", len(result))
		}
		if result[0] != 1 || result[1] != 2 || result[2] != 3 {
			t.Errorf("expected [1 2 3], got %v", result)
		}
	})

	t.Run("works with zero values", func(t *testing.T) {
		some := maybe.Just(0)
		result := some.OrElseGet(func(err error) int { return 42 })

		if result != 0 {
			t.Errorf("expected 0 (the actual value), got %d", result)
		}
	})
}

func TestSome_OrElseDefault(t *testing.T) {
	t.Run("returns the value and ignores default", func(t *testing.T) {
		some := maybe.Just(42)
		result := some.OrElseDefault(0)

		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})

	t.Run("returns string value ignoring default", func(t *testing.T) {
		some := maybe.Just("hello")
		result := some.OrElseDefault("default")

		if result != "hello" {
			t.Errorf("expected 'hello', got %s", result)
		}
	})

	t.Run("works with different types", func(t *testing.T) {
		some := maybe.Just([]int{1, 2, 3})
		result := some.OrElseDefault([]int{})

		if len(result) != 3 {
			t.Errorf("expected slice length 3, got %d", len(result))
		}
		if result[0] != 1 || result[1] != 2 || result[2] != 3 {
			t.Errorf("expected [1 2 3], got %v", result)
		}
	})

	t.Run("works with zero values", func(t *testing.T) {
		some := maybe.Just(0)
		result := some.OrElseDefault(42)

		if result != 0 {
			t.Errorf("expected 0 (the actual value), got %d", result)
		}
	})

	t.Run("can use same value as default", func(t *testing.T) {
		some := maybe.Just(10)
		result := some.OrElseDefault(10)

		if result != 10 {
			t.Errorf("expected 10, got %d", result)
		}
	})
}

func TestSome_MatchThen(t *testing.T) {
	t.Run("executes someFn and returns original Some", func(t *testing.T) {
		some := maybe.Just(42)
		var capturedValue int
		someCalled := false
		noneCalled := false
		failureCalled := false

		result := some.MatchThen(
			func(x int) {
				someCalled = true
				capturedValue = x
			},
			func() { noneCalled = true },
			func(err error) { failureCalled = true },
		)

		if !someCalled {
			t.Error("someFn should be called")
		}
		if noneCalled {
			t.Error("noneFn should not be called")
		}
		if failureCalled {
			t.Error("failureFn should not be called")
		}
		if capturedValue != 42 {
			t.Errorf("expected captured value 42, got %d", capturedValue)
		}

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("MatchThen should return Some type")
		}
		value, _ := resultSome.Get()
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})

	t.Run("can be used for logging", func(t *testing.T) {
		some := maybe.Just("hello")
		var log string

		result := some.MatchThen(
			func(x string) { log = "Got value: " + x },
			func() { log = "No value" },
			func(err error) { log = "Error: " + err.Error() },
		)

		if log != "Got value: hello" {
			t.Errorf("expected 'Got value: hello', got %s", log)
		}

		resultSome, ok := result.(maybe.Some[string])
		if !ok {
			t.Fatal("MatchThen should return Some type")
		}
		value, _ := resultSome.Get()
		if value != "hello" {
			t.Errorf("expected 'hello', got %s", value)
		}
	})

	t.Run("catches panic in someFn and converts to Failure", func(t *testing.T) {
		some := maybe.Just(10)

		result := some.MatchThen(
			func(x int) { panic("someFn panic") },
			func() {},
			func(err error) {},
		)

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MatchThen should return Failure when someFn panics")
		}
		_, err := failure.Get()
		if err.Error() != "someFn panic" {
			t.Errorf("expected panic message, got %s", err.Error())
		}
	})

	t.Run("catches panic with error type in someFn", func(t *testing.T) {
		some := maybe.Just(5)
		testErr := errors.New("test error")

		result := some.MatchThen(
			func(x int) { panic(testErr) },
			func() {},
			func(err error) {},
		)

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("MatchThen should return Failure when someFn panics")
		}
		_, err := failure.Get()
		if err != testErr {
			t.Errorf("expected %v, got %v", testErr, err)
		}
	})

	t.Run("noneFn and failureFn can panic but never called", func(t *testing.T) {
		some := maybe.Just(100)
		someCalled := false

		result := some.MatchThen(
			func(x int) { someCalled = true },
			func() { panic("noneFn should not be called") },
			func(err error) { panic("failureFn should not be called") },
		)

		if !someCalled {
			t.Error("someFn should be called")
		}

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("MatchThen should return Some type")
		}
		value, _ := resultSome.Get()
		if value != 100 {
			t.Errorf("expected 100, got %d", value)
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		var sideEffect string

		result := maybe.Just(5).
			MatchThen(
				func(x int) { sideEffect = fmt.Sprintf("Processing %d", x) },
				func() { sideEffect = "none" },
				func(err error) { sideEffect = "error" },
			).
			Map(func(x int) int { return x * 2 })

		if sideEffect != "Processing 5" {
			t.Errorf("expected 'Processing 5', got %s", sideEffect)
		}

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("chained operations should return Some")
		}
		value, _ := resultSome.Get()
		if value != 10 {
			t.Errorf("expected 10, got %v", value)
		}
	})

	t.Run("can be chained multiple times", func(t *testing.T) {
		var log []string

		result := maybe.Just(10).
			MatchThen(
				func(x int) { log = append(log, "first") },
				func() { log = append(log, "none") },
				func(err error) { log = append(log, "error") },
			).
			MatchThen(
				func(x int) { log = append(log, "second") },
				func() { log = append(log, "none") },
				func(err error) { log = append(log, "error") },
			)

		if len(log) != 2 || log[0] != "first" || log[1] != "second" {
			t.Errorf("expected [first second], got %v", log)
		}

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("chained MatchThen should return Some type")
		}
		value, _ := resultSome.Get()
		if value != 10 {
			t.Errorf("expected 10, got %d", value)
		}
	})

	t.Run("does not change value", func(t *testing.T) {
		original := maybe.Just(42)
		result := original.MatchThen(
			func(x int) { /* no-op */ },
			func() {},
			func(err error) {},
		)

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("MatchThen should return Some type")
		}
		value, _ := resultSome.Get()
		if value != 42 {
			t.Errorf("value should remain 42, got %d", value)
		}
	})
}


func TestSome_MapIfEmpty(t *testing.T) {
	t.Run("returns original Some unchanged", func(t *testing.T) {
		some := maybe.Just(42)
		called := false

		result := some.MapIfEmpty(func() (int, error) {
			called = true
			return 100, nil
		})

		if called {
			t.Error("MapIfEmpty should not call function for Some")
		}

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("MapIfEmpty should return Some for Some")
		}
		value, _ := resultSome.Get()
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})
}

func TestSome_MapIfFailed(t *testing.T) {
	t.Run("returns original Some unchanged", func(t *testing.T) {
		some := maybe.Just(42)
		called := false

		result := some.MapIfFailed(func(err error) (int, error) {
			called = true
			return 100, nil
		})

		if called {
			t.Error("MapIfFailed should not call function for Some")
		}

		resultSome, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("MapIfFailed should return Some for Some")
		}
		value, _ := resultSome.Get()
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})
}
