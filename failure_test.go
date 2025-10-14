package lwfp

import (
	"errors"
	"testing"
)

func TestFailure_GetError(t *testing.T) {
	t.Run("returns the wrapped error", func(t *testing.T) {
		err := errors.New("test error")
		failure := Failure[int]{e: err}
		if failure.GetError() != err {
			t.Errorf("expected %v, got %v", err, failure.GetError())
		}
	})

	t.Run("returns different error messages", func(t *testing.T) {
		err := errors.New("another error")
		failure := Fail[string](err)
		failureTyped := failure.(Failure[string])
		if failureTyped.GetError().Error() != "another error" {
			t.Errorf("expected 'another error', got %s", failureTyped.GetError().Error())
		}
	})
}

func TestFailure_Map(t *testing.T) {
	t.Run("propagates error and ignores function", func(t *testing.T) {
		err := errors.New("original error")
		failure := Fail[int](err)
		result := failure.Map(func(x int) any { return x * 2 })

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("Failure.Map should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		err := errors.New("test error")
		failure := Fail[int](err)
		executed := false
		result := failure.Map(func(x int) any {
			executed = true
			return x * 2
		})

		if executed {
			t.Error("Failure.Map should not execute the function")
		}

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("Failure.Map should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		err := errors.New("test error")
		failure := Fail[int](err)
		result := failure.Map(func(x int) any {
			panic("this should never be called")
		})

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("Failure.Map should return Failure type without executing function")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("preserves error through multiple Map calls", func(t *testing.T) {
		err := errors.New("persistent error")
		result := Fail[int](err).
			Map(func(x int) any { return x * 2 })

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("chained Map should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})
}

func TestFailure_FlatMap(t *testing.T) {
	t.Run("propagates error and ignores function", func(t *testing.T) {
		err := errors.New("original error")
		failure := Fail[int](err)
		result := failure.FlatMap(func(x int) Maybe[any] {
			return Just[any](x * 2)
		})

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("Failure.FlatMap should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		err := errors.New("test error")
		failure := Fail[int](err)
		executed := false
		result := failure.FlatMap(func(x int) Maybe[any] {
			executed = true
			return Just[any](x * 2)
		})

		if executed {
			t.Error("Failure.FlatMap should not execute the function")
		}

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("Failure.FlatMap should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		err := errors.New("test error")
		failure := Fail[string](err)
		result := failure.FlatMap(func(x string) Maybe[any] {
			panic("this should never be called")
		})

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("Failure.FlatMap should return Failure type without executing function")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("preserves error through multiple FlatMap calls", func(t *testing.T) {
		err := errors.New("persistent error")
		result := Fail[int](err).
			FlatMap(func(x int) Maybe[any] { return Just[any](x * 2) })

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("chained FlatMap should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("preserves error through mixed operations", func(t *testing.T) {
		err := errors.New("persistent error")
		result := Fail[int](err).
			Map(func(x int) any { return x * 2 })

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("mixed operations should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("can be used in railway-oriented programming", func(t *testing.T) {
		err := errors.New("validation failed")
		result := Fail[int](err).
			FlatMap(func(x int) Maybe[any] {
				if x > 0 {
					return Just[any](x * 2)
				}
				return Empty[any]()
			})

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("railway pattern should preserve Failure")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})
}
