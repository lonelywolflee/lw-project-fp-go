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
