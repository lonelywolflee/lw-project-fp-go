package lwfp

import (
	"errors"
	"fmt"
)

func Do[T any](fn func() Maybe[T]) (result Maybe[T]) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				result = Fail[T](err)
			} else {
				result = Fail[T](errors.New(fmt.Sprint(r)))
			}
		}
	}()

	return fn()
}
