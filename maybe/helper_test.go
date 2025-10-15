package maybe_test

import (
	"errors"
	"testing"

	"github.com/lonelywolflee/lw-project-fp-go/maybe"
)

func TestDo(t *testing.T) {
	t.Run("returns result when no panic occurs", func(t *testing.T) {
		result := maybe.Do(func() maybe.Maybe[int] {
			return maybe.Just(42)
		})

		some, ok := result.(maybe.Some[int])
		if !ok {
			t.Fatal("Do should return Some type when no panic")
		}
		if some.GetValue() != 42 {
			t.Errorf("expected 42, got %d", some.GetValue())
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
		if failure.GetError() != err {
			t.Errorf("expected %v, got %v", err, failure.GetError())
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
		if failure.GetError().Error() != "something went wrong" {
			t.Errorf("expected 'something went wrong', got %s", failure.GetError().Error())
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
		if failure.GetError() != testErr {
			t.Errorf("expected %v, got %v", testErr, failure.GetError())
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
		if failure.GetError().Error() != "123" {
			t.Errorf("expected '123', got %s", failure.GetError().Error())
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
		if some.GetValue() != "hello" {
			t.Errorf("expected 'hello', got %s", some.GetValue())
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
		if some.GetValue() != 25 {
			t.Errorf("expected 25, got %d", some.GetValue())
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
		if failure.GetError().Error() != "nested panic" {
			t.Errorf("expected 'nested panic', got %s", failure.GetError().Error())
		}
	})
}
