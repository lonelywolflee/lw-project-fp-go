package lwfp

import (
	"errors"
	"testing"
)

func TestDo(t *testing.T) {
	t.Run("returns result when no panic occurs", func(t *testing.T) {
		result := Do(func() Maybe[int] {
			return Just(42)
		})

		some, ok := result.(Some[int])
		if !ok {
			t.Fatal("Do should return Some type when no panic")
		}
		if some.GetValue() != 42 {
			t.Errorf("expected 42, got %d", some.GetValue())
		}
	})

	t.Run("returns Empty when function returns Empty", func(t *testing.T) {
		result := Do(func() Maybe[int] {
			return Empty[int]()
		})

		_, ok := result.(None[int])
		if !ok {
			t.Fatal("Do should return None type when function returns Empty")
		}
	})

	t.Run("returns Failure when function returns Failure", func(t *testing.T) {
		err := errors.New("test error")
		result := Do(func() Maybe[int] {
			return Fail[int](err)
		})

		failure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when function returns Failure")
		}
		if failure.GetError() != err {
			t.Errorf("expected %v, got %v", err, failure.GetError())
		}
	})

	t.Run("catches panic with string and converts to Failure", func(t *testing.T) {
		result := Do(func() Maybe[int] {
			panic("something went wrong")
		})

		failure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when panic occurs")
		}
		if failure.GetError().Error() != "something went wrong" {
			t.Errorf("expected 'something went wrong', got %s", failure.GetError().Error())
		}
	})

	t.Run("catches panic with error type and wraps it", func(t *testing.T) {
		testErr := errors.New("panic error")
		result := Do(func() Maybe[int] {
			panic(testErr)
		})

		failure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when panic occurs")
		}
		if failure.GetError() != testErr {
			t.Errorf("expected %v, got %v", testErr, failure.GetError())
		}
	})

	t.Run("catches panic with integer and converts to error", func(t *testing.T) {
		result := Do(func() Maybe[int] {
			panic(123)
		})

		failure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when panic with integer occurs")
		}
		if failure.GetError().Error() != "123" {
			t.Errorf("expected '123', got %s", failure.GetError().Error())
		}
	})

	t.Run("catches panic with nil pointer dereference", func(t *testing.T) {
		result := Do(func() Maybe[int] {
			var ptr *int
			_ = *ptr // This will panic
			return Just(42)
		})

		_, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when nil pointer panic occurs")
		}
	})

	t.Run("catches panic with slice out of bounds", func(t *testing.T) {
		result := Do(func() Maybe[int] {
			slice := []int{1, 2, 3}
			_ = slice[10] // This will panic
			return Just(42)
		})

		_, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when out of bounds panic occurs")
		}
	})

	t.Run("works with different return types", func(t *testing.T) {
		result := Do(func() Maybe[string] {
			return Just("hello")
		})

		some, ok := result.(Some[string])
		if !ok {
			t.Fatal("Do should return Some type")
		}
		if some.GetValue() != "hello" {
			t.Errorf("expected 'hello', got %s", some.GetValue())
		}
	})

	t.Run("handles complex operations without panic", func(t *testing.T) {
		result := Do(func() Maybe[int] {
			value := 10
			value *= 2
			value += 5
			return Just(value)
		})

		some, ok := result.(Some[int])
		if !ok {
			t.Fatal("Do should return Some type")
		}
		if some.GetValue() != 25 {
			t.Errorf("expected 25, got %d", some.GetValue())
		}
	})

	t.Run("handles panic in nested operation", func(t *testing.T) {
		result := Do(func() Maybe[int] {
			nestedFunc := func() int {
				panic("nested panic")
			}
			return Just(nestedFunc())
		})

		failure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Do should return Failure type when nested panic occurs")
		}
		if failure.GetError().Error() != "nested panic" {
			t.Errorf("expected 'nested panic', got %s", failure.GetError().Error())
		}
	})
}
