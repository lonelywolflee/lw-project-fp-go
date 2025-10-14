# lwfp - Functional Programming Toolkit for Go

[![Go Version](https://img.shields.io/badge/Go-1.25.2+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)](https://github.com/lonelywolflee/lwfp)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A lightweight, type-safe functional programming library for Go that implements the **Maybe monad** pattern. Built with Go generics for compile-time type safety and zero runtime overhead.

## Features

- **Type-Safe**: Full generic support for compile-time type checking
- **Zero Dependencies**: Pure Go implementation with no external dependencies
- **Panic Recovery**: Automatic panic-to-error conversion
- **Railway-Oriented Programming**: Clean error handling with railway pattern
- **100% Test Coverage**: Thoroughly tested with comprehensive test suite
- **Well Documented**: Complete API documentation with examples

## Installation

```bash
go get github.com/lonelywolflee/lwfp
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/lonelywolflee/lwfp"
)

func main() {
    // Create a Maybe with a value
    result := lwfp.Just(10).
        Map(func(x int) any { return x * 2 }).
        Map(func(x int) any { return x + 5 })

    // Extract the value
    if some, ok := result.(lwfp.Some[any]); ok {
        fmt.Println(some.GetValue()) // Output: 25
    }
}
```

## Core Concepts

### Maybe Monad

The `Maybe[T]` type represents an optional value with three possible states:

| Type | Description | Use Case |
|------|-------------|----------|
| `Some[T]` | Contains a value | Successful computation |
| `None[T]` | Represents absence of value | Empty result without error |
| `Failure[T]` | Contains an error | Failed computation |

### Creating Maybe Values

```go
// Create a Some (value present)
value := lwfp.Just(42)

// Create a None (no value)
empty := lwfp.Empty[int]()

// Create a Failure (error state)
failed := lwfp.Fail[int](errors.New("something went wrong"))
```

## Usage Examples

### Basic Transformations

```go
// Map transforms the value inside Maybe
result := lwfp.Just(5).
    Map(func(x int) any { return x * 2 })

// result contains Some(10)
```

### Chaining Operations

```go
// FlatMap prevents nested Maybe structures
result := lwfp.Just(5).
    FlatMap(func(x int) lwfp.Maybe[any] {
        if x > 0 {
            return lwfp.Just[any](x * 2)
        }
        return lwfp.Empty[any]()
    })
```

### Error Handling

```go
// Automatic panic recovery
result := lwfp.Just(10).
    Map(func(x int) any {
        if x > 5 {
            panic("value too large")
        }
        return x * 2
    })

// result is a Failure containing the error
if failure, ok := result.(lwfp.Failure[any]); ok {
    fmt.Println(failure.GetError()) // "value too large"
}
```

### Railway-Oriented Programming

```go
func validateAge(age int) lwfp.Maybe[any] {
    if age < 0 {
        return lwfp.Fail[any](errors.New("age cannot be negative"))
    }
    if age > 150 {
        return lwfp.Fail[any](errors.New("age too high"))
    }
    return lwfp.Just[any](age)
}

result := validateAge(25).
    Map(func(age int) any { return age * 2 }).
    Map(func(doubled int) any { return doubled + 10 })

// Error propagates through the chain
// Success path only executes if all steps succeed
```

### Using the Do Helper

```go
// Do catches panics and converts them to Failures
result := lwfp.Do(func() lwfp.Maybe[int] {
    // Risky operation that might panic
    value := riskyOperation()
    return lwfp.Just(value)
})
```

## API Reference

### Types

#### `Maybe[T]` Interface
```go
type Maybe[T any] interface {
    Map(fn func(T) any) Maybe[any]
    FlatMap(fn func(T) Maybe[any]) Maybe[any]
}
```

#### `Some[T]` Struct
```go
type Some[T any] struct { /* ... */ }

func (s Some[T]) GetValue() T
func (s Some[T]) Map(fn func(T) any) Maybe[any]
func (s Some[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any]
```

#### `None[T]` Struct
```go
type None[T any] struct{}

func (n None[T]) Get() any // returns nil
func (n None[T]) Map(fn func(T) any) Maybe[any]
func (n None[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any]
```

#### `Failure[T]` Struct
```go
type Failure[T any] struct { /* ... */ }

func (f Failure[T]) GetError() error
func (f Failure[T]) Map(fn func(T) any) Maybe[any]
func (f Failure[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any]
```

### Constructor Functions

| Function | Description |
|----------|-------------|
| `Just[T](v T) Maybe[T]` | Creates a Maybe containing a value |
| `Empty[T]() Maybe[T]` | Creates an empty Maybe (None) |
| `Fail[T](e error) Maybe[T]` | Creates a Maybe containing an error |

### Helper Functions

| Function | Description |
|----------|-------------|
| `Do[T](fn func() Maybe[T]) Maybe[T]` | Executes a function with panic recovery |

## Pattern Matching Example

```go
func handleResult(m lwfp.Maybe[int]) string {
    switch v := m.(type) {
    case lwfp.Some[int]:
        return fmt.Sprintf("Got value: %d", v.GetValue())
    case lwfp.None[int]:
        return "No value"
    case lwfp.Failure[int]:
        return fmt.Sprintf("Error: %s", v.GetError())
    default:
        return "Unknown state"
    }
}
```

## Testing

Run tests with coverage:

```bash
# Run all tests
go test -v

# Run with coverage
go test -cover

# Generate coverage report
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

1. Clone the repository
```bash
git clone https://github.com/lonelywolflee/lwfp.git
cd lwfp
```

2. Run tests
```bash
go test -v
```

3. Check coverage
```bash
go test -cover
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Inspiration

This project was created out of necessity and inspired by:
- Scala's `Option` type
- Haskell's `Maybe` monad
- Rust's `Result` and `Option` types
- Railway-oriented programming pattern

## Author

**LonelyWolfLee**

## Acknowledgments

- Thanks to the Go team for adding generics support
- Inspired by functional programming principles from various languages
- Built with the railway-oriented programming pattern in mind

---

**Note**: This library requires Go 1.18 or higher for generics support.
