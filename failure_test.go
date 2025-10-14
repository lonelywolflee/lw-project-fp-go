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

func TestFailure_Filter(t *testing.T) {
	t.Run("propagates error and ignores predicate", func(t *testing.T) {
		err := errors.New("original error")
		failure := Fail[int](err)
		result := failure.Filter(func(x int) bool { return x > 5 })

		resultFailure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Failure.Filter should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("does not execute the predicate function", func(t *testing.T) {
		err := errors.New("test error")
		failure := Fail[int](err)
		executed := false
		result := failure.Filter(func(x int) bool {
			executed = true
			return true
		})

		if executed {
			t.Error("Failure.Filter should not execute the predicate function")
		}

		resultFailure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Failure.Filter should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("predicate can panic but never called", func(t *testing.T) {
		err := errors.New("test error")
		failure := Fail[string](err)
		result := failure.Filter(func(x string) bool {
			panic("this should never be called")
		})

		resultFailure, ok := result.(Failure[string])
		if !ok {
			t.Fatal("Failure.Filter should return Failure type without executing predicate")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		err := errors.New("persistent error")
		result := Fail[int](err).
			Filter(func(x int) bool { return x > 5 }).
			Map(func(x int) any { return x * 2 })

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("chained Filter and Map should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("preserves error in railway pattern", func(t *testing.T) {
		err := errors.New("validation failed")
		result := Fail[int](err).
			Filter(func(x int) bool { return x > 0 })

		resultFailure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Filter should preserve Failure in railway pattern")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})
}

func TestFailure_Then(t *testing.T) {
	t.Run("propagates error and ignores function", func(t *testing.T) {
		err := errors.New("original error")
		failure := Fail[int](err)
		result := failure.Then(func(x int) { /* no-op */ })

		resultFailure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Failure.Then should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		err := errors.New("test error")
		failure := Fail[int](err)
		executed := false
		result := failure.Then(func(x int) {
			executed = true
		})

		if executed {
			t.Error("Failure.Then should not execute the function")
		}

		resultFailure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Failure.Then should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		err := errors.New("test error")
		failure := Fail[string](err)
		result := failure.Then(func(x string) {
			panic("this should never be called")
		})

		resultFailure, ok := result.(Failure[string])
		if !ok {
			t.Fatal("Failure.Then should return Failure type without executing function")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		err := errors.New("persistent error")
		var sideEffect int

		result := Fail[int](err).
			Then(func(x int) { sideEffect = x }).
			Map(func(x int) any { return x * 2 })

		if sideEffect != 0 {
			t.Errorf("side effect should not be triggered, got %d", sideEffect)
		}

		resultFailure, ok := result.(Failure[any])
		if !ok {
			t.Fatal("chained Then and Map should return Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("preserves error through multiple Then calls", func(t *testing.T) {
		err := errors.New("validation failed")
		var callCount int

		result := Fail[int](err).
			Then(func(x int) { callCount++ }).
			Then(func(x int) { callCount++ }).
			Then(func(x int) { callCount++ })

		if callCount != 0 {
			t.Errorf("no Then calls should execute, got %d calls", callCount)
		}

		resultFailure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("multiple Then calls should preserve Failure type")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})

	t.Run("preserves error in railway pattern with Then", func(t *testing.T) {
		err := errors.New("validation failed")
		result := Fail[int](err).
			Then(func(x int) { /* log */ })

		resultFailure, ok := result.(Failure[int])
		if !ok {
			t.Fatal("Then should preserve Failure in railway pattern")
		}
		if resultFailure.GetError() != err {
			t.Errorf("expected %v, got %v", err, resultFailure.GetError())
		}
	})
}
