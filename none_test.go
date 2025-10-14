package lwfp

import (
	"testing"
)

func TestNone_Get(t *testing.T) {
	t.Run("returns nil", func(t *testing.T) {
		none := Empty[int]()
		noneTyped := none.(None[int])
		if noneTyped.Get() != nil {
			t.Error("None.Get() should return nil")
		}
	})

	t.Run("returns nil for different types", func(t *testing.T) {
		none := Empty[string]()
		noneTyped := none.(None[string])
		if noneTyped.Get() != nil {
			t.Error("None.Get() should return nil")
		}
	})
}

func TestNone_Map(t *testing.T) {
	t.Run("returns Empty and ignores function", func(t *testing.T) {
		none := Empty[int]()
		result := none.Map(func(x int) any { return x * 2 })

		_, ok := result.(None[any])
		if !ok {
			t.Fatal("None.Map should return None type")
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		none := Empty[int]()
		executed := false
		result := none.Map(func(x int) any {
			executed = true
			return x * 2
		})

		if executed {
			t.Error("None.Map should not execute the function")
		}

		_, ok := result.(None[any])
		if !ok {
			t.Fatal("None.Map should return None type")
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		none := Empty[int]()
		result := none.Map(func(x int) any {
			panic("this should never be called")
		})

		_, ok := result.(None[any])
		if !ok {
			t.Fatal("None.Map should return None type without executing function")
		}
	})
}

func TestNone_FlatMap(t *testing.T) {
	t.Run("returns Empty and ignores function", func(t *testing.T) {
		none := Empty[int]()
		result := none.FlatMap(func(x int) Maybe[any] {
			return Just[any](x * 2)
		})

		_, ok := result.(None[any])
		if !ok {
			t.Fatal("None.FlatMap should return None type")
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		none := Empty[int]()
		executed := false
		result := none.FlatMap(func(x int) Maybe[any] {
			executed = true
			return Just[any](x * 2)
		})

		if executed {
			t.Error("None.FlatMap should not execute the function")
		}

		_, ok := result.(None[any])
		if !ok {
			t.Fatal("None.FlatMap should return None type")
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		none := Empty[string]()
		result := none.FlatMap(func(x string) Maybe[any] {
			panic("this should never be called")
		})

		_, ok := result.(None[any])
		if !ok {
			t.Fatal("None.FlatMap should return None type without executing function")
		}
	})

	t.Run("chains with Map operation", func(t *testing.T) {
		result := Empty[int]().
			Map(func(x int) any { return x * 2 })

		_, ok := result.(None[any])
		if !ok {
			t.Fatal("chained operations on None should return None type")
		}
	})
}

func TestNone_Filter(t *testing.T) {
	t.Run("returns None and ignores predicate", func(t *testing.T) {
		none := Empty[int]()
		result := none.Filter(func(x int) bool { return x > 5 })

		_, ok := result.(None[int])
		if !ok {
			t.Fatal("None.Filter should return None type")
		}
	})

	t.Run("does not execute the predicate function", func(t *testing.T) {
		none := Empty[int]()
		executed := false
		result := none.Filter(func(x int) bool {
			executed = true
			return true
		})

		if executed {
			t.Error("None.Filter should not execute the predicate function")
		}

		_, ok := result.(None[int])
		if !ok {
			t.Fatal("None.Filter should return None type")
		}
	})

	t.Run("predicate can panic but never called", func(t *testing.T) {
		none := Empty[string]()
		result := none.Filter(func(x string) bool {
			panic("this should never be called")
		})

		_, ok := result.(None[string])
		if !ok {
			t.Fatal("None.Filter should return None type without executing predicate")
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		result := Empty[int]().
			Filter(func(x int) bool { return x > 5 }).
			Map(func(x int) any { return x * 2 })

		_, ok := result.(None[any])
		if !ok {
			t.Fatal("chained Filter and Map on None should return None type")
		}
	})

	t.Run("preserves None through Filter", func(t *testing.T) {
		result := Empty[int]().
			Filter(func(x int) bool { return x > 10 })

		_, ok := result.(None[int])
		if !ok {
			t.Fatal("Filter on None should preserve None type")
		}
	})
}

func TestNone_Then(t *testing.T) {
	t.Run("returns None and ignores function", func(t *testing.T) {
		none := Empty[int]()
		result := none.Then(func(x int) { /* no-op */ })

		_, ok := result.(None[int])
		if !ok {
			t.Fatal("None.Then should return None type")
		}
	})

	t.Run("does not execute the function", func(t *testing.T) {
		none := Empty[int]()
		executed := false
		result := none.Then(func(x int) {
			executed = true
		})

		if executed {
			t.Error("None.Then should not execute the function")
		}

		_, ok := result.(None[int])
		if !ok {
			t.Fatal("None.Then should return None type")
		}
	})

	t.Run("function can panic but never called", func(t *testing.T) {
		none := Empty[string]()
		result := none.Then(func(x string) {
			panic("this should never be called")
		})

		_, ok := result.(None[string])
		if !ok {
			t.Fatal("None.Then should return None type without executing function")
		}
	})

	t.Run("can be chained with Map", func(t *testing.T) {
		var sideEffect int

		result := Empty[int]().
			Then(func(x int) { sideEffect = x }).
			Map(func(x int) any { return x * 2 })

		if sideEffect != 0 {
			t.Errorf("side effect should not be triggered, got %d", sideEffect)
		}

		_, ok := result.(None[any])
		if !ok {
			t.Fatal("chained Then and Map on None should return None type")
		}
	})

	t.Run("preserves None through multiple Then calls", func(t *testing.T) {
		var callCount int

		result := Empty[int]().
			Then(func(x int) { callCount++ }).
			Then(func(x int) { callCount++ }).
			Then(func(x int) { callCount++ })

		if callCount != 0 {
			t.Errorf("no Then calls should execute, got %d calls", callCount)
		}

		_, ok := result.(None[int])
		if !ok {
			t.Fatal("multiple Then calls on None should preserve None type")
		}
	})
}
