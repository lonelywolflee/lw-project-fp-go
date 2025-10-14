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
go get github.com/lonelywolflee/lw-project-fp-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/lonelywolflee/lw-project-fp-go"
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

### Filtering Values

```go
// Filter keeps values that satisfy a predicate
result := lwfp.Just(10).
    Filter(func(x int) bool { return x > 5 })
// result contains Some(10)

result := lwfp.Just(3).
    Filter(func(x int) bool { return x > 5 })
// result contains None

// Filter can be chained with other operations
result := lwfp.Just(10).
    Filter(func(x int) bool { return x > 5 }).
    Map(func(x int) any { return x * 2 })
// result contains Some(20)
```

### Side Effects with Then

```go
// Then executes a function for side effects without changing the value
result := lwfp.Just(10).
    Then(func(x int) { fmt.Printf("Processing: %d\n", x) }).
    Map(func(x int) any { return x * 2 })
// Prints "Processing: 10", result contains Some(20)

// Useful for logging in processing pipelines
result := lwfp.Just("data").
    Then(func(x string) { log.Info("Received", x) }).
    Filter(func(x string) bool { return len(x) > 0 }).
    Then(func(x string) { log.Info("Validated", x) }).
    Map(func(x string) any { return strings.ToUpper(x) })
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

### Extracting Values with Defaults

```go
// OrElseGet provides a computed default value
value := lwfp.Just(42).OrElseGet(func() int { return 0 })
// Returns: 42 (the actual value)

value := lwfp.Empty[int]().OrElseGet(func() int { return 0 })
// Returns: 0 (computed default)

value := lwfp.Fail[int](err).OrElseGet(func() int {
    return computeDefault()
})
// Returns: result of computeDefault()

// OrElseDefault provides a constant default value
value := lwfp.Just(42).OrElseDefault(0)
// Returns: 42 (the actual value)

value := lwfp.Empty[int]().OrElseDefault(0)
// Returns: 0 (constant default)

value := lwfp.Fail[int](err).OrElseDefault(0)
// Returns: 0 (constant default)

// Practical example: config with fallback
config := fetchConfig().
    Filter(func(c Config) bool { return c.IsValid() }).
    OrElseDefault(DefaultConfig)
```

### Pattern Matching with MatchThen

```go
// MatchThen provides exhaustive pattern matching for all Maybe states
result := fetchUser(id).MatchThen(
    func(user User) {
        log.Info("Found user", user.Name)
    },
    func() {
        log.Warn("User not found")
    },
    func(err error) {
        log.Error("Database error", err)
    },
)

// Returns the original Maybe unchanged, allowing for chaining
value := processData().
    MatchThen(
        func(data string) { metrics.RecordSuccess() },
        func() { metrics.RecordEmpty() },
        func(err error) { metrics.RecordError(err) },
    ).
    Map(func(data string) any { return transform(data) }).
    OrElseDefault(defaultData)

// Practical example: HTTP request handling
response := makeAPICall().
    MatchThen(
        func(data Response) {
            fmt.Printf("Success: %d items\n", len(data.Items))
        },
        func() {
            fmt.Println("No data received")
        },
        func(err error) {
            fmt.Printf("API Error: %v\n", err)
        },
    )
```

## API Reference

### Types

#### `Maybe[T]` Interface
```go
type Maybe[T any] interface {
    Map(fn func(T) any) Maybe[any]
    FlatMap(fn func(T) Maybe[any]) Maybe[any]
    Filter(fn func(T) bool) Maybe[T]
    Then(fn func(T)) Maybe[T]
    OrElseGet(fn func() T) T
    OrElseDefault(v T) T
    MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T]
}
```

#### `Some[T]` Struct
```go
type Some[T any] struct { /* ... */ }

func (s Some[T]) GetValue() T
func (s Some[T]) Map(fn func(T) any) Maybe[any]
func (s Some[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any]
func (s Some[T]) Filter(fn func(T) bool) Maybe[T]
func (s Some[T]) Then(fn func(T)) Maybe[T]
func (s Some[T]) OrElseGet(fn func() T) T
func (s Some[T]) OrElseDefault(v T) T
func (s Some[T]) MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T]
```

#### `None[T]` Struct
```go
type None[T any] struct{}

func (n None[T]) Get() any // returns nil
func (n None[T]) Map(fn func(T) any) Maybe[any]
func (n None[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any]
func (n None[T]) Filter(fn func(T) bool) Maybe[T]
func (n None[T]) Then(fn func(T)) Maybe[T]
func (n None[T]) OrElseGet(fn func() T) T
func (n None[T]) OrElseDefault(v T) T
func (n None[T]) MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T]
```

#### `Failure[T]` Struct
```go
type Failure[T any] struct { /* ... */ }

func (f Failure[T]) GetError() error
func (f Failure[T]) Map(fn func(T) any) Maybe[any]
func (f Failure[T]) FlatMap(fn func(T) Maybe[any]) Maybe[any]
func (f Failure[T]) Filter(fn func(T) bool) Maybe[T]
func (f Failure[T]) Then(fn func(T)) Maybe[T]
func (f Failure[T]) OrElseGet(fn func() T) T
func (f Failure[T]) OrElseDefault(v T) T
func (f Failure[T]) MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T]
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
git clone https://github.com/lonelywolflee/lw-project-fp-go.git
cd lw-project-fp-go
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
