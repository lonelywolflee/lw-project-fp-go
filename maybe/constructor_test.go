package maybe_test

import (
	"errors"
	"testing"

	"github.com/lonelywolflee/lw-project-fp-go/maybe"
)

func TestJust(t *testing.T) {
	t.Run("creates Some with integer value", func(t *testing.T) {
		some := maybe.Just(42)
		value, _ := some.Get()
		if value != 42 {
			t.Errorf("expected 42, got %d", value)
		}
	})

	t.Run("creates Some with string value", func(t *testing.T) {
		some := maybe.Just("hello")
		value, _ := some.Get()
		if value != "hello" {
			t.Errorf("expected 'hello', got %s", value)
		}
	})

	t.Run("creates Some with nil pointer", func(t *testing.T) {
		var ptr *int = nil
		some := maybe.Just(ptr)
		value, _ := some.Get()
		if value != nil {
			t.Error("expected nil pointer")
		}
	})
}

func TestEmpty(t *testing.T) {
	t.Run("creates None for int type", func(t *testing.T) {
		none := maybe.Empty[int]()
		// Just verify it's created, type is enforced by generics
		_ = none
	})

	t.Run("creates None for string type", func(t *testing.T) {
		none := maybe.Empty[string]()
		// Just verify it's created, type is enforced by generics
		_ = none
	})
}

func TestFail(t *testing.T) {
	t.Run("creates Failure with error", func(t *testing.T) {
		err := errors.New("test error")
		failure := maybe.Failed[int](err)
		_, gotErr := failure.Get()
		if gotErr != err {
			t.Errorf("expected %v, got %v", err, gotErr)
		}
	})

	t.Run("creates Failure with different error message", func(t *testing.T) {
		err := errors.New("another error")
		failure := maybe.Failed[string](err)
		_, gotErr := failure.Get()
		if gotErr.Error() != "another error" {
			t.Errorf("expected 'another error', got %s", gotErr.Error())
		}
	})
}
