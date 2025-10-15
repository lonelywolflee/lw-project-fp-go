package maybe_test

import (
	"errors"
	"testing"

	"github.com/lonelywolflee/lw-project-fp-go/maybe"
)

// customMaybe is a test helper type that implements Maybe[int] interface
// but is not actually Some, None, or Failure. This is used to test the
// default case in helper functions Map and FlatMap.
type customMaybe struct{}

func (customMaybe) Map(fn func(int) int) maybe.Maybe[int] {
	return maybe.Empty[int]()
}
func (customMaybe) FlatMap(fn func(int) maybe.Maybe[int]) maybe.Maybe[int] {
	return maybe.Empty[int]()
}
func (customMaybe) Filter(fn func(int) bool) maybe.Maybe[int] {
	return maybe.Empty[int]()
}
func (customMaybe) Then(fn func(int)) maybe.Maybe[int] {
	return maybe.Empty[int]()
}
func (customMaybe) Get() (int, error) {
	return 0, nil // Acts like None but isn't actually None type
}
func (customMaybe) OrElseGet(fn func(error) int) int {
	return 0
}
func (customMaybe) OrElseDefault(v int) int {
	return v
}
func (customMaybe) FailIfEmpty(fn func() error) maybe.Maybe[int] {
	return maybe.Empty[int]()
}
func (customMaybe) MatchThen(someFn func(int), noneFn func(), failureFn func(error)) maybe.Maybe[int] {
	return maybe.Empty[int]()
}

func TestDo(t *testing.T) {
	t.Run("returns result when no panic occurs", func(t *testing.T) {
		result := maybe.Do(func() maybe.Maybe[int] {
			return maybe.Just(42)
		})

		some, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("Do should return Some type when no panic")
		}
		value, _ := some.Get()
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})

	t.Run("returns Empty when function returns Empty", func(t *testing.T) {
		result := maybe.Do(func() maybe.Maybe[int] {
			return maybe.Empty[int]()
		})

		_, ok := result.(maybe.None[int])
		if !ok {
			t.Fatal("Do should return None type when function returns Empty")
		}
	})

	t.Run("returns Failure when function returns Failure", func(t *testing.T) {
		err := errors.New("test error")
		result := maybe.Do(func() maybe.Maybe[int] {
			return maybe.Fail[int](err)
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when function returns Failure")
		}
		_, gotErr := failure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("catches panic with string and converts to Failure", func(t *testing.T) {
		result := maybe.Do(func() maybe.Maybe[int] {
			panic("something went wrong")
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when panic occurs")
		}
		_, err := failure.Get()
		if err.Error() != "something went wrong" {
			t.Errorf("expected 'something went wrong', got %s", err.Error())
		}
	})

	t.Run("catches panic with error type and wraps it", func(t *testing.T) {
		testErr := errors.New("panic error")
		result := maybe.Do(func() maybe.Maybe[int] {
			panic(testErr)
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when panic occurs")
		}
		_, err := failure.Get()
		if err != testErr {
			t.Errorf("expected %v, got %v", testErr, err)
		}
	})

	t.Run("catches panic with integer and converts to error", func(t *testing.T) {
		result := maybe.Do(func() maybe.Maybe[int] {
			panic(123)
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when panic with integer occurs")
		}
		_, err := failure.Get()
		if err.Error() != "123" {
			t.Errorf("expected '123', got %s", err.Error())
		}
	})

	t.Run("catches panic with nil pointer dereference", func(t *testing.T) {
		result := maybe.Do(func() maybe.Maybe[int] {
			var ptr *int
			_ = *ptr // This will panic
			return maybe.Just(42)
		})

		_, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when nil pointer panic occurs")
		}
	})

	t.Run("catches panic with slice out of bounds", func(t *testing.T) {
		result := maybe.Do(func() maybe.Maybe[int] {
			slice := []int{1, 2, 3}
			_ = slice[10] // This will panic
			return maybe.Just(42)
		})

		_, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when out of bounds panic occurs")
		}
	})

	t.Run("works with different return types", func(t *testing.T) {
		result := maybe.Do(func() maybe.Maybe[string] {
			return maybe.Just("hello")
		})

		some, ok := result.(maybe.Some[string])
		if !ok {
			t.Fatal("Do should return Some type")
		}
		value, _ := some.Get()
		if value != "hello" {
			t.Errorf("expected 'hello', got %s", value)
		}
	})

	t.Run("handles complex operations without panic", func(t *testing.T) {
		result := maybe.Do(func() maybe.Maybe[int] {
			value := 10
			value *= 2
			value += 5
			return maybe.Just(value)
		})

		some, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("Do should return Some type")
		}
		value, _ := some.Get()
		if value != 25 {
			t.Errorf("expected 25, got %d", value)
		}
	})

	t.Run("handles panic in nested operation", func(t *testing.T) {
		result := maybe.Do(func() maybe.Maybe[int] {
			nestedFunc := func() int {
				panic("nested panic")
			}
			return maybe.Just(nestedFunc())
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when nested panic occurs")
		}
		_, err := failure.Get()
		if err.Error() != "nested panic" {
			t.Errorf("expected 'nested panic', got %s", err.Error())
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("transforms Some value to different type", func(t *testing.T) {
		// int to string
		result := maybe.Map(maybe.Just(42), func(x int) string {
			return "value: " + string(rune(x))
		})

		some, ok := result.(maybe.Some[string])
		if !ok {
			t.Fatal("Map should return Some[string]")
		}
		value, _ := some.Get()
		if value != "value: *" {
			t.Errorf("expected 'value: *', got %s", value)
		}
	})

	t.Run("transforms None to different type", func(t *testing.T) {
		result := maybe.Map(maybe.Empty[int](), func(x int) string {
			return "value"
		})

		_, ok := result.(maybe.None[string])
		if !ok {
			t.Fatal("Map should return None[string] for None[int]")
		}
	})

	t.Run("propagates Failure to different type", func(t *testing.T) {
		originalErr := errors.New("original error")
		result := maybe.Map(maybe.Fail[int](originalErr), func(x int) string {
			return "value"
		})

		failure, ok := result.(maybe.Failure[string])
		if !ok {
			t.Fatal("Map should return Failure[string] for Failure[int]")
		}
		_, err := failure.Get()
		if err != originalErr {
			t.Errorf("expected original error, got %v", err)
		}
	})

	t.Run("catches panic and converts to Failure", func(t *testing.T) {
		result := maybe.Map(maybe.Just(42), func(x int) string {
			panic("panic in map function")
		})

		failure, ok := result.(maybe.Failure[string])
		if !ok {
			t.Fatal("Map should return Failure when function panics")
		}
		_, err := failure.Get()
		if err.Error() != "panic in map function" {
			t.Errorf("expected 'panic in map function', got %s", err.Error())
		}
	})

	t.Run("converts int to string using strconv", func(t *testing.T) {
		result := maybe.Map(maybe.Just(123), func(x int) string {
			return string(rune(x + '0'))
		})

		some, ok := result.(maybe.Some[string])
		if !ok {
			t.Fatal("Map should return Some[string]")
		}
		value, _ := some.Get()
		if len(value) == 0 {
			t.Error("expected non-empty string")
		}
	})

	t.Run("converts string to int length", func(t *testing.T) {
		result := maybe.Map(maybe.Just("hello"), func(s string) int {
			return len(s)
		})

		some, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("Map should return Some[int]")
		}
		value, _ := some.Get()
		if value != 5 {
			t.Errorf("expected 5, got %d", value)
		}
	})

	t.Run("can be chained with method calls", func(t *testing.T) {
		result := maybe.Map(
			maybe.Just(10).Filter(func(x int) bool { return x > 5 }),
			func(x int) string { return "passed" },
		)

		some, ok := result.(maybe.Some[string])
		if !ok {
			t.Fatal("Map should return Some[string]")
		}
		value, _ := some.Get()
		if value != "passed" {
			t.Errorf("expected 'passed', got %s", value)
		}
	})

	t.Run("handles zero values correctly", func(t *testing.T) {
		result := maybe.Map(maybe.Just(0), func(x int) bool {
			return x == 0
		})

		some, ok := result.(maybe.Some[bool])
		if !ok {
			t.Fatal("Map should return Some[bool]")
		}
		value, _ := some.Get()
		if !value {
			t.Error("expected true")
		}
	})

	t.Run("handles unknown Maybe implementation (default case)", func(t *testing.T) {
		// Use the custom Maybe implementation to test the default case
		custom := customMaybe{}
		result := maybe.Map(custom, func(x int) string {
			return "should not reach"
		})

		// Should return Empty[string]() due to default case
		_, ok := result.(maybe.None[string])
		if !ok {
			t.Fatal("Map should return None[string] for unknown Maybe implementation")
		}
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("transforms Some value to different type with Maybe", func(t *testing.T) {
		result := maybe.FlatMap(maybe.Just(42), func(x int) maybe.Maybe[string] {
			if x > 0 {
				return maybe.Just("positive")
			}
			return maybe.Empty[string]()
		})

		some, ok := result.(maybe.Some[string])
		if !ok {
			t.Fatal("FlatMap should return Some[string]")
		}
		value, _ := some.Get()
		if value != "positive" {
			t.Errorf("expected 'positive', got %s", value)
		}
	})

	t.Run("transforms Some to None based on condition", func(t *testing.T) {
		result := maybe.FlatMap(maybe.Just(-5), func(x int) maybe.Maybe[string] {
			if x > 0 {
				return maybe.Just("positive")
			}
			return maybe.Empty[string]()
		})

		_, ok := result.(maybe.None[string])
		if !ok {
			t.Fatal("FlatMap should return None[string] when function returns Empty")
		}
	})

	t.Run("transforms None to different type", func(t *testing.T) {
		result := maybe.FlatMap(maybe.Empty[int](), func(x int) maybe.Maybe[string] {
			return maybe.Just("value")
		})

		_, ok := result.(maybe.None[string])
		if !ok {
			t.Fatal("FlatMap should return None[string] for None[int]")
		}
	})

	t.Run("propagates Failure to different type", func(t *testing.T) {
		originalErr := errors.New("original error")
		result := maybe.FlatMap(maybe.Fail[int](originalErr), func(x int) maybe.Maybe[string] {
			return maybe.Just("value")
		})

		failure, ok := result.(maybe.Failure[string])
		if !ok {
			t.Fatal("FlatMap should return Failure[string] for Failure[int]")
		}
		_, err := failure.Get()
		if err != originalErr {
			t.Errorf("expected original error, got %v", err)
		}
	})

	t.Run("catches panic and converts to Failure", func(t *testing.T) {
		result := maybe.FlatMap(maybe.Just(42), func(x int) maybe.Maybe[string] {
			panic("panic in flatmap function")
		})

		failure, ok := result.(maybe.Failure[string])
		if !ok {
			t.Fatal("FlatMap should return Failure when function panics")
		}
		_, err := failure.Get()
		if err.Error() != "panic in flatmap function" {
			t.Errorf("expected 'panic in flatmap function', got %s", err.Error())
		}
	})

	t.Run("transforms to Failure based on validation", func(t *testing.T) {
		result := maybe.FlatMap(maybe.Just("invalid"), func(s string) maybe.Maybe[int] {
			if s == "valid" {
				return maybe.Just(100)
			}
			return maybe.Fail[int](errors.New("validation failed"))
		})

		failure, ok := result.(maybe.Failure[int])
		if !ok {
			t.Fatal("FlatMap should return Failure when validation fails")
		}
		_, err := failure.Get()
		if err.Error() != "validation failed" {
			t.Errorf("expected 'validation failed', got %s", err.Error())
		}
	})

	t.Run("can chain multiple transformations", func(t *testing.T) {
		// First transformation: string to int length
		step1 := maybe.FlatMap(maybe.Just("hello"), func(s string) maybe.Maybe[int] {
			return maybe.Just(len(s))
		})

		// Second transformation: int to bool
		step2 := maybe.FlatMap(step1, func(x int) maybe.Maybe[bool] {
			return maybe.Just(x > 3)
		})

		some, ok := step2.(maybe.Some[bool])
		if !ok {
			t.Fatal("FlatMap should return Some[bool]")
		}
		value, _ := some.Get()
		if !value {
			t.Error("expected true")
		}
	})

	t.Run("flattens nested Maybe structures", func(t *testing.T) {
		// Without FlatMap, this would be Maybe[Maybe[string]]
		result := maybe.FlatMap(maybe.Just(5), func(x int) maybe.Maybe[string] {
			if x > 0 {
				return maybe.Just("value")
			}
			return maybe.Empty[string]()
		})

		some, ok := result.(maybe.Some[string])
		if !ok {
			t.Fatal("FlatMap should flatten and return Some[string], not Maybe[Maybe[string]]")
		}
		value, _ := some.Get()
		if value != "value" {
			t.Errorf("expected 'value', got %s", value)
		}
	})

	t.Run("can be chained with Filter", func(t *testing.T) {
		result := maybe.FlatMap(
			maybe.Just(10).Filter(func(x int) bool { return x > 5 }),
			func(x int) maybe.Maybe[string] {
				return maybe.Just("valid")
			},
		)

		some, ok := result.(maybe.Some[string])
		if !ok {
			t.Fatal("FlatMap should return Some[string]")
		}
		value, _ := some.Get()
		if value != "valid" {
			t.Errorf("expected 'valid', got %s", value)
		}
	})

	t.Run("handles Filter returning None", func(t *testing.T) {
		result := maybe.FlatMap(
			maybe.Just(3).Filter(func(x int) bool { return x > 5 }),
			func(x int) maybe.Maybe[string] {
				return maybe.Just("valid")
			},
		)

		_, ok := result.(maybe.None[string])
		if !ok {
			t.Fatal("FlatMap should return None[string] when Filter returns None")
		}
	})

	t.Run("handles unknown Maybe implementation (default case)", func(t *testing.T) {
		// Use the custom Maybe implementation to test the default case
		custom := customMaybe{}
		result := maybe.FlatMap(custom, func(x int) maybe.Maybe[string] {
			return maybe.Just("should not reach")
		})

		// Should return Empty[string]() due to default case
		_, ok := result.(maybe.None[string])
		if !ok {
			t.Fatal("FlatMap should return None[string] for unknown Maybe implementation")
		}
	})
}
