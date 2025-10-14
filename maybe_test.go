package lwfp

import (
	"errors"
	"testing"
)

func TestJust(t *testing.T) {
	t.Run("creates Some with integer value", func(t *testing.T) {
		maybe := Just(42)
		some, ok := maybe.(Some[int])
		if !ok {
			t.Fatal("Just should return Some type")
		}
		if some.GetValue() != 42 {
			t.Errorf("expected 42, got %d", some.GetValue())
		}
	})

	t.Run("creates Some with string value", func(t *testing.T) {
		maybe := Just("hello")
		some, ok := maybe.(Some[string])
		if !ok {
			t.Fatal("Just should return Some type")
		}
		if some.GetValue() != "hello" {
			t.Errorf("expected 'hello', got %s", some.GetValue())
		}
	})

	t.Run("creates Some with nil pointer", func(t *testing.T) {
		var ptr *int = nil
		maybe := Just(ptr)
		some, ok := maybe.(Some[*int])
		if !ok {
			t.Fatal("Just should return Some type")
		}
		if some.GetValue() != nil {
			t.Error("expected nil pointer")
		}
	})
}

func TestEmpty(t *testing.T) {
	t.Run("creates None for int type", func(t *testing.T) {
		maybe := Empty[int]()
		_, ok := maybe.(None[int])
		if !ok {
			t.Fatal("Empty should return None type")
		}
	})

	t.Run("creates None for string type", func(t *testing.T) {
		maybe := Empty[string]()
		_, ok := maybe.(None[string])
		if !ok {
			t.Fatal("Empty should return None type")
		}
	})
}

func TestFail(t *testing.T) {
	t.Run("creates Failure with error", func(t *testing.T) {
		err := errors.New("test error")
		maybe := Fail[int](err)
		failure, ok := maybe.(Failure[int])
		if !ok {
			t.Fatal("Fail should return Failure type")
		}
		if failure.GetError() != err {
			t.Errorf("expected %v, got %v", err, failure.GetError())
		}
	})

	t.Run("creates Failure with different error message", func(t *testing.T) {
		err := errors.New("another error")
		maybe := Fail[string](err)
		failure, ok := maybe.(Failure[string])
		if !ok {
			t.Fatal("Fail should return Failure type")
		}
		if failure.GetError().Error() != "another error" {
			t.Errorf("expected 'another error', got %s", failure.GetError().Error())
		}
	})
}
