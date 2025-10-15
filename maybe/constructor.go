package maybe

// Just creates a Maybe that contains a value (Some).
// Use this when you have a valid value to wrap.
//
// Example:
//
//	maybe := Just(42)
//	// To get the value with type safety, cast to Some and use GetValue():
//	value := maybe.(Some[int]).GetValue() // returns 42 as int
func Just[T any](v T) Some[T] {
	return Some[T]{v: v}
}

// Empty creates a Maybe that represents the absence of a value (None).
// Use this when you want to explicitly represent "no value" without an error.
//
// Example:
//
//	maybe := Empty[int]()
//	// Check if it's None:
//	if _, ok := maybe.(None[int]); ok {
//	    // Handle empty case
//	}
func Empty[T any]() None[T] {
	return None[T]{}
}

// Fail creates a Maybe that represents an error state (Failure).
// Use this when you want to wrap an error in the Maybe monad.
//
// Example:
//
//	maybe := Fail[int](errors.New("something went wrong"))
//	// To get the error, cast to Failure and use GetError():
//	err := maybe.(Failure[int]).GetError() // returns the error
func Fail[T any](e error) Failure[T] {
	return Failure[T]{e: e}
}
