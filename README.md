# lwfp - Functional Programming Toolkit for Go

[![Go Version](https://img.shields.io/badge/Go-1.18+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)](https://github.com/lonelywolflee/lwfp)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A lightweight, type-safe functional programming library for Go that implements the **Maybe monad** pattern. Built with Go generics for compile-time type safety and zero runtime overhead. Designed to embrace Go's type system constraints while providing powerful functional programming abstractions.

## Features

- **Type-Safe**: Full generic support for compile-time type checking
- **Zero Dependencies**: Pure Go implementation with no external dependencies
- **Go Interop**: Seamless integration with Go's standard `(T, error)` pattern via `ToMaybe` and `Try`
- **Panic Recovery**: Automatic panic-to-error conversion
- **Railway-Oriented Programming**: Clean error handling with railway pattern
- **Go-Idiomatic Design**: Works with Go's type system, not against it
- **100% Test Coverage**: Thoroughly tested with comprehensive test suite
- **Well Documented**: Complete API documentation with examples

## Installation

```bash
go get github.com/lonelywolflee/lw-project-fp-go
```

## Design Philosophy & Type System

### Why Type-Constrained Map/FlatMap?

Go's type system has specific constraints that influenced this library's design. Understanding these constraints helps you write better, more idiomatic code.

#### The Go Generics Limitation

```go
// ❌ This is NOT possible in Go - methods cannot have their own type parameters
type Maybe[T any] interface {
    Map[R any](fn func(T) R) Maybe[R]  // Compiler error!
}

// ✅ This is what Go allows - same-type transformations
type Maybe[T any] interface {
    Map(fn func(T) T) Maybe[T]  // Same type only
}
```

**Why this matters:**
- Go methods cannot introduce new type parameters beyond the receiver's type
- This is a fundamental language design decision, not a limitation
- We embrace this constraint rather than fight it

#### Our Solution: Two-Tier API

**1. Method Chaining (Same Type)**
```go
// For transformations that keep the same type
result := maybe.Just(10).
    Map(func(x int) int { return x * 2 }).     // int → int
    Filter(func(x int) bool { return x > 5 }). // keeps int
    Map(func(x int) int { return x + 10 })     // int → int
// result: Just(30)
```

**2. Helper Functions (Type Conversion)**
```go
// For transformations that change types
import "strconv"

result := maybe.Map(maybe.Just(42), strconv.Itoa)
// int → string: Just("42")

// Or with inline function
result := maybe.Map(maybe.Just(10), func(x int) string {
    return fmt.Sprintf("value: %d", x)
})
// result: Just("value: 10")
```

### Design Principles

1. **Explicitness Over Magic**: Type conversions are clearly visible in the code
2. **Go-Native Patterns**: Follows Go idioms rather than forcing functional patterns
3. **Compile-Time Safety**: All type errors caught at compile time
4. **Railway-Oriented Programming**: Errors propagate cleanly through chains
5. **Practical Pragmatism**: Most operations don't need type conversion

### When to Use Each Approach

**Use Method Chaining (`.Map()`, `.FlatMap()`) when:**
- ✅ Transforming within the same type (int → int, string → string)
- ✅ Building processing pipelines with Filter, Then, MatchThen
- ✅ You want fluent, readable chains
- ✅ Most of your application logic

**Use Helper Functions (`Map()`, `FlatMap()`) when:**
- ✅ Converting between types (int → string, string → int)
- ✅ Integrating with external APIs that return different types
- ✅ Parsing or serialization operations
- ✅ Type conversion is the primary goal

## Quick Start

### 1. Basic Same-Type Transformations

The most common pattern - transforming values while keeping the same type:

```go
package main

import (
    "fmt"
    "github.com/lonelywolflee/lw-project-fp-go/maybe"
)

func main() {
    // Chain operations on the same type
    result := maybe.Just(10).
        Map(func(x int) int { return x * 2 }).      // 10 → 20
        Filter(func(x int) bool { return x > 15 }). // keeps 20
        Map(func(x int) int { return x + 5 })       // 20 → 25

    // Extract the value
    value, ok, err := result.Get()
    if ok && err == nil {
        fmt.Println(value) // Output: 25
    }
}
```

### 2. Type Conversion with Helper Functions

When you need to convert between types, use the helper functions:

```go
package main

import (
    "fmt"
    "strconv"
    "github.com/lonelywolflee/lw-project-fp-go/maybe"
)

func main() {
    // Convert int to string using helper function
    result := maybe.Map(maybe.Just(42), strconv.Itoa)

    value, _, _ := result.Get()
    fmt.Println(value) // Output: "42"

    // Can also use inline functions
    result2 := maybe.Map(maybe.Just(100), func(x int) string {
        return fmt.Sprintf("Score: %d", x)
    })

    value2, _, _ := result2.Get()
    fmt.Println(value2) // Output: "Score: 100"
}
```

### 3. Interop with Go Standard Library (ToMaybe & Try)

Easily integrate Go's standard `(T, error)` pattern with Maybe:

```go
package main

import (
    "fmt"
    "strconv"
    "os"
    "github.com/lonelywolflee/lw-project-fp-go/maybe"
)

func main() {
    // ToMaybe: Convert (value, error) tuple to Maybe
    result := maybe.ToMaybe(strconv.Atoi("42"))
    fmt.Println(result) // Just(42)

    result2 := maybe.ToMaybe(strconv.Atoi("invalid"))
    fmt.Println(result2) // Failed(error)

    // Chain with other operations
    parsed := maybe.ToMaybe(strconv.Atoi("100")).
        Filter(func(x int) bool { return x > 0 }).
        Map(func(x int) int { return x * 2 })

    value, _, _ := parsed.Get()
    fmt.Println(value) // Output: 200

    // Try: Execute function with both error AND panic handling
    config := maybe.Try(func() ([]byte, error) {
        return os.ReadFile("config.json")
    }).Filter(func(data []byte) bool {
        return len(data) > 0
    }).MapIfEmpty(func() ([]byte, error) {
        return nil, fmt.Errorf("config file is empty")
    })

    // Try catches both errors and panics
    safeParse := maybe.Try(func() (int, error) {
        return strconv.Atoi("42")
    }).Map(func(x int) int {
        return x * 2
    })

    value2, _, _ := safeParse.Get()
    fmt.Println(value2) // Output: 84
}
```

**When to use:**
- **ToMaybe**: When you already have a `(T, error)` tuple from standard library functions
- **Try**: When you need panic protection in addition to error handling

### 4. Error Handling with Railway-Oriented Programming

Build pipelines where errors automatically propagate:

```go
package main

import (
    "errors"
    "fmt"
    "github.com/lonelywolflee/lw-project-fp-go/maybe"
)

func validateAge(age int) maybe.Maybe[int] {
    if age < 0 {
        return maybe.Failed[int](errors.New("age cannot be negative"))
    }
    if age > 150 {
        return maybe.Failed[int](errors.New("age too high"))
    }
    return maybe.Just(age)
}

func main() {
    // Success case - all operations execute
    result := validateAge(25).
        Map(func(age int) int { return age * 2 }). // Executes
        Map(func(age int) int { return age + 10 }) // Executes

    value, _, _ := result.Get()
    fmt.Println(value) // Output: 60

    // Error case - operations are skipped after error
    result2 := validateAge(-5).
        Map(func(age int) int { return age * 2 }). // Skipped
        Map(func(age int) int { return age + 10 }) // Skipped

    _, _, err := result2.Get()
    fmt.Println(err) // Output: age cannot be negative
}
```

### 5. Practical Example: Data Validation Pipeline

A real-world example combining multiple features:

```go
package main

import (
    "errors"
    "fmt"
    "strings"
    "github.com/lonelywolflee/lw-project-fp-go/maybe"
)

type User struct {
    Name  string
    Email string
    Age   int
}

func validateUser(name, email string, age int) maybe.Maybe[User] {
    // Validate name using method chaining
    validatedName := maybe.Just(name).
        Filter(func(n string) bool { return len(strings.TrimSpace(n)) > 0 }).
        MapIfEmpty(func() (string, error) { return "", errors.New("name is required") }).
        Filter(func(n string) bool { return len(n) < 100 }).
        MapIfEmpty(func() (string, error) { return "", errors.New("name too long") })

    // Transform to User struct using helper FlatMap (type conversion: string → User)
    return maybe.FlatMap(validatedName, func(n string) maybe.Maybe[User] {
        // Validate email
        if !strings.Contains(email, "@") {
            return maybe.Failed[User](errors.New("invalid email"))
        }

        // Validate age
        if age < 0 || age > 150 {
            return maybe.Failed[User](errors.New("invalid age"))
        }

        return maybe.Just(User{
            Name:  n,
            Email: email,
            Age:   age,
        })
    })
}

func main() {
    // Valid user
    result := validateUser("John Doe", "john@example.com", 30)
    if user, ok, err := result.Get(); ok && err == nil {
        fmt.Printf("Valid user: %+v\n", user)
    }

    // Invalid user - empty name
    result2 := validateUser("", "john@example.com", 30)
    if _, _, err := result2.Get(); err != nil {
        fmt.Printf("Error: %v\n", err) // Output: Error: name is required
    }

    // Invalid user - bad email
    result3 := validateUser("John", "invalid-email", 30)
    if _, _, err := result3.Get(); err != nil {
        fmt.Printf("Error: %v\n", err) // Output: Error: invalid email
    }
}
```

### 6. Combining Type Conversion with Validation

Use helper functions with method chaining for powerful pipelines:

```go
package main

import (
    "fmt"
    "strconv"
    "github.com/lonelywolflee/lw-project-fp-go/maybe"
)

func parseAndValidate(input string) maybe.Maybe[string] {
    // Parse string to int using ToMaybe (converts (value, error) to Maybe)
    parsedInt := maybe.ToMaybe(strconv.Atoi(input))

    // Validate using method chaining (same type: int → int)
    validatedInt := parsedInt.
        Filter(func(x int) bool { return x > 0 }).
        MapIfEmpty(func() (int, error) { return 0, fmt.Errorf("value must be positive") })

    // Format to string using helper Map (type conversion: int → string)
    return maybe.Map(validatedInt, func(x int) string {
        return fmt.Sprintf("Valid: %d", x)
    })
}

func main() {
    // Valid input
    result := parseAndValidate("42")
    fmt.Println(result.OrElseDefault("error")) // Output: Valid: 42

    // Invalid input - not a number
    result2 := parseAndValidate("abc")
    _, _, err := result2.Get()
    fmt.Println(err) // Output: strconv.Atoi: parsing "abc": invalid syntax

    // Invalid input - negative
    result3 := parseAndValidate("-5")
    _, _, err2 := result3.Get()
    fmt.Println(err2) // Output: value must be positive
}
```

### 6. Using `Then` for Side Effects

Log or track progress without changing values:

```go
package main

import (
    "fmt"
    "log"
    "strings"
    "github.com/lonelywolflee/lw-project-fp-go/maybe"
)

func processData(data string) maybe.Maybe[string] {
    return maybe.Just(data).
        Then(func(s string) { log.Printf("Received: %s", s) }).
        Filter(func(s string) bool { return len(s) > 0 }).
        Then(func(s string) { log.Printf("Validated: %s", s) }).
        Map(func(s string) string { return strings.ToUpper(s) }).
        Then(func(s string) { log.Printf("Processed: %s", s) })
}

func main() {
    result := processData("hello")
    value, _, _ := result.Get()
    fmt.Println(value) // Output: HELLO
    // Logs:
    // Received: hello
    // Validated: hello
    // Processed: HELLO
}
```

### 7. Pattern Matching with MatchThen

Handle all possible states explicitly:

```go
package main

import (
    "fmt"
    "github.com/lonelywolflee/lw-project-fp-go/maybe"
)

func processResult(m maybe.Maybe[int]) {
    m.MatchThen(
        func(value int) {
            fmt.Printf("Success: Got value %d\n", value)
        },
        func() {
            fmt.Println("Empty: No value present")
        },
        func(err error) {
            fmt.Printf("Error: %v\n", err)
        },
    )
}

func main() {
    processResult(maybe.Just(42))              // Success: Got value 42
    processResult(maybe.Empty[int]())          // Empty: No value present
    processResult(maybe.Failed[int](fmt.Errorf("failed"))) // Error: failed
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
value := maybe.Just(42)

// Create a None (no value)
empty := maybe.Empty[int]()

// Create a Failure (error state)
failed := maybe.Failed[int](errors.New("something went wrong"))
```

## Usage Examples

### Basic Same-Type Transformations

```go
// Map transforms the value inside Maybe (same type)
result := maybe.Just(5).
    Map(func(x int) int { return x * 2 })

// result contains Some(10)
```

### Type Conversion with Helper Functions

```go
// For converting between types, use helper functions
import "strconv"

// Convert int to string
result := maybe.Map(maybe.Just(42), strconv.Itoa)
// result contains Some("42")

// Convert with custom logic
result := maybe.Map(maybe.Just(5), func(x int) string {
    return fmt.Sprintf("value: %d", x)
})
// result contains Some("value: 5")
```

### Chaining Same-Type Operations

```go
// FlatMap prevents nested Maybe structures (same type)
result := maybe.Just(5).
    FlatMap(func(x int) maybe.Maybe[int] {
        if x > 0 {
            return maybe.Just(x * 2)
        }
        return maybe.Empty[int]()
    })
// result contains Some(10)
```

### Error Handling with Automatic Panic Recovery

```go
// Automatic panic recovery
result := maybe.Just(10).
    Map(func(x int) int {
        if x > 5 {
            panic("value too large")
        }
        return x * 2
    })

// result is a Failure containing the error
_, _, err := result.Get()
if err != nil {
    fmt.Println(err) // "value too large"
}
```

### Railway-Oriented Programming

```go
func validateAge(age int) maybe.Maybe[int] {
    if age < 0 {
        return maybe.Failed[int](errors.New("age cannot be negative"))
    }
    if age > 150 {
        return maybe.Failed[int](errors.New("age too high"))
    }
    return maybe.Just(age)
}

result := validateAge(25).
    Map(func(age int) int { return age * 2 }).
    Map(func(doubled int) int { return doubled + 10 })

// Error propagates through the chain
// Success path only executes if all steps succeed
```

### Filtering Values

```go
// Filter keeps values that satisfy a predicate
result := maybe.Just(10).
    Filter(func(x int) bool { return x > 5 })
// result contains Some(10)

result := maybe.Just(3).
    Filter(func(x int) bool { return x > 5 })
// result contains None

// Filter can be chained with other operations (same type)
result := maybe.Just(10).
    Filter(func(x int) bool { return x > 5 }).
    Map(func(x int) int { return x * 2 })
// result contains Some(20)
```

### Side Effects with Then

```go
// Then executes a function for side effects without changing the value
result := maybe.Just(10).
    Then(func(x int) { fmt.Printf("Processing: %d\n", x) }).
    Map(func(x int) int { return x * 2 })
// Prints "Processing: 10", result contains Some(20)

// Useful for logging in processing pipelines
result := maybe.Just("data").
    Then(func(x string) { log.Info("Received", x) }).
    Filter(func(x string) bool { return len(x) > 0 }).
    Then(func(x string) { log.Info("Validated", x) }).
    Map(func(x string) string { return strings.ToUpper(x) })
```

### Using the Do Helper

```go
// Do catches panics and converts them to Failures
result := maybe.Do(func() maybe.Maybe[int] {
    // Risky operation that might panic
    value := riskyOperation()
    return maybe.Just(value)
})
```

### Extracting Values with Get

```go
// Get provides a Go-idiomatic way to extract values with presence checking and error handling
// Returns: (value, ok, error)
// - ok=true: value is present (Some state)
// - ok=false: no value (None or Failure state)
value, ok, err := maybe.Just(42).Get()
if ok && err == nil {
    fmt.Println(value) // 42
}

// For None: returns zero value, false, and nil error
value, ok, err := maybe.Empty[int]().Get()
// value = 0, ok = false, err = nil

// For Failure: returns zero value, false, and the error
value, ok, err := maybe.Failed[int](errors.New("failed")).Get()
// value = 0, ok = false, err = error("failed")

// Practical example: database query
func findUser(id int) (User, bool, error) {
    result := queryDatabase(id) // returns Maybe[User]
    return result.Get()
}

// Or if you only care about the value and error:
func findUserSimple(id int) (User, error) {
    result := queryDatabase(id) // returns Maybe[User]
    user, _, err := result.Get()
    return user, err
}
```

### Extracting Values with Defaults

```go
// OrElseGet provides a computed default value with error awareness
// Function receives: nil for None, actual error for Failure
value := maybe.Just(42).OrElseGet(func(err error) int { return 0 })
// Returns: 42 (the actual value, function not called)

value := maybe.Empty[int]().OrElseGet(func(err error) int { return 0 })
// Returns: 0 (err is nil for None)

value := maybe.Failed[int](err).OrElseGet(func(e error) int {
    log.Printf("Error occurred: %v", e)
    return computeDefault()
})
// Returns: result of computeDefault() (logs the error)

// Practical example: error-aware default computation
result := validateInput(data).OrElseGet(func(err error) int {
    if err != nil {
        log.Printf("Validation failed: %v", err)
        return -1 // Error indicator
    }
    return 0 // Empty indicator
})

// OrElseDefault provides a constant default value
value := maybe.Just(42).OrElseDefault(0)
// Returns: 42 (the actual value)

value := maybe.Empty[int]().OrElseDefault(0)
// Returns: 0 (constant default)

value := maybe.Failed[int](err).OrElseDefault(0)
// Returns: 0 (constant default)

// Practical example: config with fallback
config := fetchConfig().
    Filter(func(c Config) bool { return c.IsValid() }).
    OrElseDefault(DefaultConfig)
```

### Error Recovery and Transformation with MapIfEmpty and MapIfFailed

These methods provide dual-purpose functionality: **recovery** (converting to Some) and **error transformation** (converting to/modifying Failure).

**Basic patterns:**

```go
// MapIfEmpty Pattern 1: Recovery - convert None to Some
result := maybe.Empty[int]().MapIfEmpty(func() (int, error) {
    return 42, nil  // Provide default value
}) // Just(42)

// MapIfEmpty Pattern 2: Error Transformation - convert None to Failure
result := maybe.Empty[int]().MapIfEmpty(func() (int, error) {
    return 0, errors.New("value required")
}) // Failed[int]("value required")

// MapIfFailed Pattern 1: Recovery - convert Failure to Some
result := maybe.Failed[int](errors.New("not found")).MapIfFailed(func(err error) (int, error) {
    if errors.Is(err, ErrNotFound) {
        return 0, nil  // Recover from specific error
    }
    return 0, err  // Propagate other errors
}) // Just(0)

// MapIfFailed Pattern 2: Error Transformation - wrap or enrich errors
result := maybe.Failed[int](dbErr).MapIfFailed(func(err error) (int, error) {
    return 0, fmt.Errorf("user service error: %w", err)
}) // Failed[int](wrapped error)

// Some remains unchanged for both
result := maybe.Just(10).MapIfEmpty(func() (int, error) {
    return 42, nil  // Never called
}) // Just(10)

result := maybe.Just(10).MapIfFailed(func(err error) (int, error) {
    return 42, nil  // Never called
}) // Just(10)
```

**Practical examples:**

```go
// Example 1: Fallback chain with MapIfEmpty
config := loadPrimaryConfig().
    MapIfEmpty(func() (Config, error) {
        log.Info("Primary config not found, loading backup")
        return loadBackupConfig()
    }).
    MapIfEmpty(func() (Config, error) {
        log.Info("Backup config not found, using defaults")
        return DefaultConfig, nil
    })

// Example 2: Retry logic with MapIfFailed
data := fetchFromAPI().
    MapIfFailed(func(err error) (Data, error) {
        log.Printf("API failed: %v, trying cache", err)
        return fetchFromCache()
    }).
    MapIfFailed(func(err error) (Data, error) {
        log.Printf("Cache failed: %v, using stale data", err)
        return fetchStaleData()
    })

// Example 3: Error-specific recovery
user := getUserByID(id).
    MapIfFailed(func(err error) (User, error) {
        if errors.Is(err, ErrNotFound) {
            // Create default user for not found
            return AnonymousUser, nil
        }
        if errors.Is(err, ErrPermission) {
            // Log security event but don't recover
            log.Security("Permission denied", id)
            return User{}, err
        }
        // Propagate unexpected errors
        return User{}, err
    })

// Example 4: Combining both recovery methods
result := processInput(input).
    MapIfEmpty(func() (Result, error) {
        // Provide default when input processing returns nothing
        return DefaultResult, nil
    }).
    MapIfFailed(func(err error) (Result, error) {
        // Try to recover from processing errors
        if errors.Is(err, ErrInvalidFormat) {
            return sanitizeAndRetry(input)
        }
        return Result{}, err
    })

// Example 5: Error transformation - Converting None to domain-specific error
user := findUserInCache(id).
    MapIfEmpty(func() (User, error) {
        // Convert None to domain-specific error
        return User{}, fmt.Errorf("user %d not in cache", id)
    })

// Example 6: Error transformation - enriching errors with context
result := databaseQuery(sql).
    MapIfFailed(func(err error) (Data, error) {
        // Add context to database errors
        return Data{}, fmt.Errorf("query failed for table 'users': %w", err)
    })

// Example 7: Multi-layer error transformation
payment := processPayment(order).
    MapIfFailed(func(err error) (Payment, error) {
        // Layer 1: Convert infrastructure errors to domain errors
        if errors.Is(err, sql.ErrNoRows) {
            return Payment{}, ErrPaymentNotFound
        }
        if errors.Is(err, context.DeadlineExceeded) {
            return Payment{}, ErrPaymentTimeout
        }
        return Payment{}, fmt.Errorf("payment processing error: %w", err)
    })
```

**When to use:**

**MapIfEmpty:**
- **Recovery**: Provide default values or fallback logic for empty states
- **Error Transformation**: Convert None to Failure with custom error

**MapIfFailed:**
- **Recovery**: Convert Failure to Some by providing fallback values or retry logic
- **Error Transformation**: Wrap, enrich, or convert errors (e.g., DB errors → domain errors)

**Both:** Can be chained together for comprehensive error handling, recovery, and transformation

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
    // Same-type transformations
    Map(fn func(T) T) Maybe[T]
    FlatMap(fn func(T) Maybe[T]) Maybe[T]

    // Filtering and side effects
    Filter(fn func(T) bool) Maybe[T]
    Then(fn func(T)) Maybe[T]

    // Value extraction
    Get() (T, bool, error)
    OrElseGet(fn func(error) T) T
    OrElseDefault(v T) T

    // Error handling and recovery
    MapIfEmpty(fn func() (T, error)) Maybe[T]
    MapIfFailed(fn func(error) (T, error)) Maybe[T]
    MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T]
}
```

#### `Some[T]` Struct
```go
type Some[T any] struct { /* ... */ }

func (s Some[T]) Map(fn func(T) T) Maybe[T]
func (s Some[T]) FlatMap(fn func(T) Maybe[T]) Maybe[T]
func (s Some[T]) Filter(fn func(T) bool) Maybe[T]
func (s Some[T]) Then(fn func(T)) Maybe[T]
func (s Some[T]) Get() (T, bool, error)
func (s Some[T]) OrElseGet(fn func(error) T) T
func (s Some[T]) OrElseDefault(v T) T
func (s Some[T]) MapIfEmpty(fn func() (T, error)) Maybe[T]
func (s Some[T]) MapIfFailed(fn func(error) (T, error)) Maybe[T]
func (s Some[T]) MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T]
```

#### `None[T]` Struct
```go
type None[T any] struct{}

func (n None[T]) Map(fn func(T) T) Maybe[T]
func (n None[T]) FlatMap(fn func(T) Maybe[T]) Maybe[T]
func (n None[T]) Filter(fn func(T) bool) Maybe[T]
func (n None[T]) Then(fn func(T)) Maybe[T]
func (n None[T]) Get() (T, bool, error)
func (n None[T]) OrElseGet(fn func(error) T) T
func (n None[T]) OrElseDefault(v T) T
func (n None[T]) MapIfEmpty(fn func() (T, error)) Maybe[T]
func (n None[T]) MapIfFailed(fn func(error) (T, error)) Maybe[T]
func (n None[T]) MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T]
```

#### `Failure[T]` Struct
```go
type Failure[T any] struct { /* ... */ }

func (f Failure[T]) Map(fn func(T) T) Maybe[T]
func (f Failure[T]) FlatMap(fn func(T) Maybe[T]) Maybe[T]
func (f Failure[T]) Filter(fn func(T) bool) Maybe[T]
func (f Failure[T]) Then(fn func(T)) Maybe[T]
func (f Failure[T]) Get() (T, bool, error)
func (f Failure[T]) OrElseGet(fn func(error) T) T
func (f Failure[T]) OrElseDefault(v T) T
func (f Failure[T]) MapIfEmpty(fn func() (T, error)) Maybe[T]
func (f Failure[T]) MapIfFailed(fn func(error) (T, error)) Maybe[T]
func (f Failure[T]) MatchThen(someFn func(T), noneFn func(), failureFn func(error)) Maybe[T]
```

### Constructor Functions

| Function | Description |
|----------|-------------|
| `Just[T](v T) Some[T]` | Creates a Some containing a value |
| `Empty[T]() None[T]` | Creates an empty None |
| `Failed[T](e error) Failure[T]` | Creates a Failure containing an error |

### Helper Functions

| Function | Description |
|----------|-------------|
| `ToMaybe[T](v T, err error) Maybe[T]` | Converts Go's standard (value, error) pattern to Maybe[T] |
| `Try[T](fn func() (T, error)) Maybe[T]` | Executes a function with both error and panic handling |
| `Do[T](fn func() Maybe[T]) Maybe[T]` | Executes a function with panic recovery |
| `Map[T, R](m Maybe[T], fn func(T) R) Maybe[R]` | Transforms Maybe[T] to Maybe[R] (type conversion) |
| `FlatMap[T, R](m Maybe[T], fn func(T) Maybe[R]) Maybe[R]` | FlatMaps Maybe[T] to Maybe[R] (type conversion) |

**Key Features:**
- **ToMaybe** and **Try**: Bridge the gap between Go's standard error handling and the Maybe monad
- **Map** and **FlatMap**: Enable type conversion across different types (not possible with methods)
- **Do**: Provides panic recovery for risky operations

## Pattern Matching Example

```go
func handleResult(m maybe.Maybe[int]) string {
    value, ok, err := m.Get()
    if err != nil {
        return fmt.Sprintf("Error: %s", err)
    }

    switch m.(type) {
    case maybe.Some[int]:
        return fmt.Sprintf("Got value: %d", value)
    case maybe.None[int]:
        return "No value"
    default:
        return "Unknown state"
    }
}

// Or use the ok flag for simpler logic:
func handleResultSimple(m maybe.Maybe[int]) string {
    value, ok, err := m.Get()
    if err != nil {
        return fmt.Sprintf("Error: %s", err)
    }
    if ok {
        return fmt.Sprintf("Got value: %d", value)
    }
    return "No value"
}
```

## Testing

Run tests with coverage:

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
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
go test -v ./...
```

3. Check coverage
```bash
go test -cover ./...
```

### Project Structure

The `maybe` package is organized as follows:

- **maybe.go** - Core `Maybe[T]` interface definition (same-type transformations)
- **constructor.go** - Constructor functions (`Just`, `Empty`, `Fail`)
- **some.go** - `Some[T]` implementation (value present)
- **none.go** - `None[T]` implementation (value absent)
- **failure.go** - `Failure[T]` implementation (error state)
- **helper.go** - Helper functions (`Do` for panic recovery, `Map`/`FlatMap` for type conversion)
- **\*_test.go** - Comprehensive test suite with 100% coverage

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

## FAQ

### Why can't I use `.Map()` to convert between types?

Go's type system doesn't allow methods to introduce new type parameters. This is a language design decision, not a bug. Use the helper functions `Map[T, R]()` and `FlatMap[T, R]()` for type conversions.

### When should I use method chaining vs helper functions?

- **Method chaining** (`.Map()`, `.FlatMap()`): Use for same-type transformations (int → int, string → string)
- **Helper functions** (`Map()`, `FlatMap()`): Use for type conversions (int → string, string → bool)

### Is this library production-ready?

Yes! It has:
- 100% test coverage
- Zero dependencies
- Comprehensive documentation
- Battle-tested design patterns

### How does this compare to Option/Result types in other languages?

It's similar to:
- Rust's `Option<T>` and `Result<T, E>` combined
- Scala's `Option[T]` with error handling
- Haskell's `Maybe` monad with `Either` for errors

The key difference is adapting to Go's type system constraints.

---

**Note**: This library requires Go 1.18 or higher for generics support.
