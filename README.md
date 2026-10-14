# Go Mapper

A powerful, type-safe Go library for mapping between structs with custom transformation functions.

## Features

- **Type Safety**: Full compile-time type checking with generics
- **Performance Optimized**: Minimal overhead with direct function calls
- **Thread Safe**: Concurrent-safe mapping registry with sync.RWMutex
- **Flexible**: Support for value, pointer and interface types
- **List Mapping**: Built-in support for mapping slices of structs
- **Zero Dependencies**: Pure Go implementation with no external dependencies

## Installation

```bash
go get github.com/mdnishat27/go-mapper
```

Requires Go 1.21 or later. The package name is `mapper`.

## Quick Start

```go
package main

import (
    "fmt"

    mapper "github.com/mdnishat27/go-mapper"
)

type User struct {
    FirstName string
    LastName  string
    Age       int
}

type UserDTO struct {
    FullName string
    Age      int
    IsAdult  bool
}

func MapUserToDTO(user User) (UserDTO, error) {
    return UserDTO{
        FullName: user.FirstName + " " + user.LastName,
        Age:      user.Age,
        IsAdult:  user.Age >= 18,
    }, nil
}

func main() {
    // Register the mapping function
    mapper.RegisterMapper[User, UserDTO](MapUserToDTO)

    // Use the mapping
    user := User{FirstName: "John", LastName: "Doe", Age: 25}
    dto, err := mapper.Map[User, UserDTO](user)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Mapped DTO: %+v\n", dto)
}
```

## API Reference

### RegisterMapper[Source, Destination](fn MappingFunction[Source, Destination])

Registers a custom mapping function for the given type pair. Registering the same pair again replaces the previous function. Panics if `fn` is nil.

```go
type MappingFunction[Source any, Destination any] func(Source) (Destination, error)

// Example
mapper.RegisterMapper[User, UserDTO](MapUserToDTO)
```

### Map[Source, Destination](source Source) (Destination, error)

Maps a source struct to a destination struct using the registered mapping function.

```go
user := User{FirstName: "Jane", LastName: "Smith", Age: 30}
dto, err := mapper.Map[User, UserDTO](user)
```

### MapList[Source, Destination](sources []Source) ([]Destination, error)

Maps a slice of source structs to a slice of destination structs. It stops at the first element that fails and returns a nil slice with an error naming the failing index.

```go
users := []User{
    {FirstName: "Alice", LastName: "Johnson", Age: 28},
    {FirstName: "Bob", LastName: "Brown", Age: 35},
}

dtos, err := mapper.MapList[User, UserDTO](users)
```

### HasMapping[Source, Destination]() bool

Checks if a mapping function is registered for the given type pair.

```go
if mapper.HasMapping[User, UserDTO]() {
    // Mapping is available
}
```

### UnregisterMapper[Source, Destination]() bool

Removes the mapping function for the given type pair. Returns `false` if none was registered.

```go
removed := mapper.UnregisterMapper[User, UserDTO]()
```

### ClearRegistry()

Removes all registered mappings. Useful for testing.

```go
func TestMyFunction(t *testing.T) {
    mapper.ClearRegistry()
    // Register test mappings...
}
```

### ErrMapperNotFound

`Map` and `MapList` return an error wrapping `ErrMapperNotFound` when no mapping function is registered for the type pair.

```go
dto, err := mapper.Map[User, UserDTO](user)
if errors.Is(err, mapper.ErrMapperNotFound) {
    // No mapping registered for User -> UserDTO
}
```

## Advanced Usage

### Working with Pointers

Go Mapper supports both value and pointer types. Each combination is a separate type pair with its own registration:

```go
// Register mapping from struct to pointer
mapper.RegisterMapper[User, *UserDTO](MapUserToDTOPointer)

// Register mapping from pointer to struct
mapper.RegisterMapper[*User, UserDTO](MapPointerUserToDTO)

// Register mapping from pointer to pointer
mapper.RegisterMapper[*User, *UserDTO](MapPointerUserToDTOPointer)
```

### Complex Transformations

Your mapping functions can contain any logic:

```go
func MapUserToEmployeeDTO(user User) (EmployeeDTO, error) {
    salary := 50000 + (user.Age * 1000) // Complex logic
    if user.Age < 18 {
        return EmployeeDTO{}, errors.New("user must be at least 18")
    }

    return EmployeeDTO{
        Name:   user.FirstName + " " + user.LastName,
        Salary: salary,
        Level:  calculateLevel(user.Age),
    }, nil
}
```

### Error Handling

Go Mapper propagates errors from your mapping functions. `Map` returns them unchanged; `MapList` wraps them with the index of the failing element, so use `errors.Is` / `errors.As` to inspect them:

```go
func MapUserToDTO(user User) (UserDTO, error) {
    if user.FirstName == "" {
        return UserDTO{}, errors.New("first name cannot be empty")
    }
    // ... mapping logic
}

// Usage
dto, err := mapper.Map[User, UserDTO](user)
if err != nil {
    // Handle mapping error
    log.Printf("Mapping failed: %v", err)
    return
}
```

## Performance

Go Mapper is designed for high performance:

- **Direct Function Calls**: Reflection is only used to look up the type pair in the registry, never to copy fields
- **Minimal Allocations**: `Map` adds no allocations of its own; `MapList` allocates the result slice once
- **Concurrent Safe**: Lookups take a read lock, so concurrent mapping does not serialize

Run the benchmarks on your own machine:

```bash
# Library benchmarks
go test -run '^$' -bench . -benchmem

# Comparison with manual copy, reflection, jinzhu/copier and JSON round trip
cd benchmarks && go test -run '^$' -bench . -benchmem
```

The comparison benchmarks live in a separate module under [benchmarks](./benchmarks/) so that the libraries they compare against are not dependencies of go-mapper.

## Best Practices

1. **Register Early**: Register all your mapping functions during application initialization
2. **Handle Errors**: Always check and handle errors from mapping functions
3. **Use ClearRegistry() in Tests**: Ensure test isolation by clearing the registry
4. **Keep Mapping Functions Pure**: Avoid side effects in mapping functions for better testability
5. **Document Complex Logic**: Add comments for non-obvious transformation logic

## Examples

See the [example](./example/) directory for a complete working example demonstrating all features of the go-mapper library.

Runnable examples are in [example_test.go](./example_test.go), and the test suite in [mapper_test.go](./mapper_test.go) covers every library feature.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure `gofmt -l .`, `go vet ./...` and `go test -race ./...` are clean
5. Submit a pull request

## License

Apache License 2.0 - see the [LICENSE](LICENSE) file for details.