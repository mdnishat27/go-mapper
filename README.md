<div align="center">

<h1>🔁 go-mapper</h1>

<p><strong>Type-safe struct mapping for Go, powered by generics.</strong></p>

<p>Register each conversion once, then call it from anywhere.<br>
No reflection-based field copying, no struct tags, no code generation.</p>

<p>
  <a href="https://github.com/mdnishat27/go-mapper/actions/workflows/ci.yml"><img src="https://github.com/mdnishat27/go-mapper/actions/workflows/ci.yml/badge.svg?branch=main" alt="Tests"></a>
  <a href="https://github.com/mdnishat27/go-mapper/actions/workflows/lint.yml"><img src="https://github.com/mdnishat27/go-mapper/actions/workflows/lint.yml/badge.svg?branch=main" alt="Lint"></a>
  <a href="https://codecov.io/gh/mdnishat27/go-mapper"><img src="https://codecov.io/gh/mdnishat27/go-mapper/graph/badge.svg" alt="codecov"></a>
  <a href="https://pkg.go.dev/github.com/mdnishat27/go-mapper"><img src="https://pkg.go.dev/badge/github.com/mdnishat27/go-mapper.svg" alt="Go Reference"></a>
  <a href="https://github.com/mdnishat27/go-mapper/releases/latest"><img src="https://img.shields.io/github/v/release/mdnishat27/go-mapper" alt="Latest release"></a>
  <a href="./go.mod"><img src="https://img.shields.io/github/go-mod/go-version/mdnishat27/go-mapper" alt="Go Version"></a>
  <a href="./LICENSE"><img src="https://img.shields.io/badge/license-Apache%202.0-blue.svg" alt="License"></a>
</p>

<p>
  <a href="#installation">Installation</a> •
  <a href="#quick-start">Quick Start</a> •
  <a href="#api-at-a-glance">API</a> •
  <a href="#performance">Benchmarks</a> •
  <a href="https://pkg.go.dev/github.com/mdnishat27/go-mapper">Docs</a>
</p>

</div>

<br>

```go
mapper.RegisterMapper[User, UserDTO](MapUserToDTO)

dto, err := mapper.Map[User, UserDTO](user)            // one value
dtos, err := mapper.MapList[User, UserDTO](users)      // a whole slice
```

## Why go-mapper?

Converting between entities, DTOs and API models is everywhere in Go services. The usual options each have a catch:

| Approach | Catch |
| --- | --- |
| Reflection copiers (e.g. `copier`) | Match fields by name at runtime, so a renamed field silently stops being copied, and every call pays for reflection |
| JSON marshal/unmarshal round trip | Slow, allocates heavily and loses anything without a JSON tag |
| Hand-written functions called directly | Fast and explicit, but scattered across packages with no common way to find them, map slices or report errors |

go-mapper keeps the plain hand-written function, the fastest and most explicit option, and adds a typed registry around it:

- ✅ **Compile-time type safety.** Generic type parameters mean a mapping function with the wrong signature doesn't compile.
- ⚡ **Direct function calls.** Reflection is only used to find your function, never to copy fields. `Map` adds zero allocations and is about 11× faster than a reflection copy and about 50× faster than `copier` or a JSON round trip ([benchmarks](#performance)).
- 🧵 **Safe for concurrent use.** The registry can be read and written from any goroutine.
- 📦 **Slice mapping built in.** `MapList` preallocates the result and tells you which element failed.
- 🎯 **Predictable errors.** Your function's errors come back unchanged, and a missing mapping is detectable with `errors.Is`.
- 🪶 **Zero dependencies.** Standard library only.

## Installation

```bash
go get github.com/mdnishat27/go-mapper
```

Requires **Go 1.21+**. The package name is `mapper`.

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/mdnishat27/go-mapper"
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
    // 1. Register once, typically at startup
    mapper.RegisterMapper[User, UserDTO](MapUserToDTO)

    // 2. Map anywhere
    dto, err := mapper.Map[User, UserDTO](User{FirstName: "John", LastName: "Doe", Age: 25})
    if err != nil {
        panic(err)
    }

    fmt.Printf("%+v\n", dto)
    // {FullName:John Doe Age:25 IsAdult:true}
}
```

A complete runnable program is in [example/basic](./example/basic/main.go):

```bash
go run ./example/basic
```

## API at a Glance

| Function | Description |
| --- | --- |
| `RegisterMapper[S, D](fn)` | Registers `fn` for the `S → D` pair. Re-registering replaces it. Panics if `fn` is nil. |
| `Map[S, D](src)` | Maps one value using the registered function. |
| `MapList[S, D](srcs)` | Maps a slice. Stops at the first failure and returns an error naming the index. A nil input gives a nil output. |
| `HasMapping[S, D]()` | Reports whether a mapping is registered for the pair. |
| `UnregisterMapper[S, D]()` | Removes the mapping. Returns `false` if none was registered. |
| `ClearRegistry()` | Removes every mapping. Useful in tests. |
| `ErrMapperNotFound` | Wrapped by `Map` and `MapList` when no mapping is registered. |
| `MappingFunction[S, D]` | `func(S) (D, error)`, the signature every mapping function has. |

Full documentation with examples is on [pkg.go.dev](https://pkg.go.dev/github.com/mdnishat27/go-mapper).

## Usage

### Mapping slices

```go
users := []User{
    {FirstName: "Alice", LastName: "Johnson", Age: 28},
    {FirstName: "Bob", LastName: "Brown", Age: 35},
}

dtos, err := mapper.MapList[User, UserDTO](users)
```

### Handling errors

`Map` returns your function's error unchanged. `MapList` wraps it with the failing index, so use `errors.Is` / `errors.As` to inspect it:

```go
func MapUserToDTO(user User) (UserDTO, error) {
    if user.FirstName == "" {
        return UserDTO{}, errors.New("first name cannot be empty")
    }
    // ...
}

_, err := mapper.MapList[User, UserDTO]([]User{{FirstName: "Alice"}, {}})
fmt.Println(err) // mapping element 1: first name cannot be empty
```

Check for a missing registration with `ErrMapperNotFound`:

```go
if errors.Is(err, mapper.ErrMapperNotFound) {
    // No mapping registered for User -> UserDTO
}
```

### Pointers and interfaces

Every combination of value, pointer and interface type is its own pair with its own registration:

```go
mapper.RegisterMapper[User, *UserDTO](MapUserToDTOPointer)
mapper.RegisterMapper[*User, UserDTO](MapPointerUserToDTO)
mapper.RegisterMapper[fmt.Stringer, string](func(s fmt.Stringer) (string, error) {
    return s.String(), nil
})
```

A mapping registered for an interface type is looked up by that interface, not by the concrete type passed in.

### Managing registrations

```go
if mapper.HasMapping[User, UserDTO]() {
    // mapping is available
}

removed := mapper.UnregisterMapper[User, UserDTO]()
```

### Testing

The registry is global, so clear it at the start of tests that register mappings:

```go
func TestSomething(t *testing.T) {
    mapper.ClearRegistry()
    mapper.RegisterMapper[User, UserDTO](MapUserToDTO)
    // ...
}
```

## Performance

Your mapping function is called directly, so `Map` costs one read-locked map lookup on top of the function itself, with **zero allocations**.

Converting one struct with six fields (`int64`, two `string`s, `decimal.Decimal`, `bool`, `time.Time`):

| Approach | Time per op | Memory per op | Allocs per op | vs go-mapper |
| --- | ---: | ---: | ---: | ---: |
| Manual function call (baseline) | 8.8 ns | 0 B | 0 | — |
| **go-mapper** | **52 ns** | **0 B** | **0** | — |
| Hand-rolled reflection copy | 597 ns | 192 B | 2 | 11× slower |
| JSON marshal + unmarshal | 2,587 ns | 720 B | 16 | 49× slower |
| [`jinzhu/copier`](https://github.com/jinzhu/copier) | 2,942 ns | 664 B | 21 | 56× slower |

`MapList` over 1,000 elements takes about 6.4 µs (≈6.4 ns per element) with a single allocation for the result slice, because the registry lookup happens once per call, not once per element.

<sub>Median of 3 runs. Intel Core i7-1185G7, Windows, amd64, Go 1.25. Your numbers will vary by machine.</sub>

To compare go-mapper with a manual copy, a hand-rolled reflection copy, [`jinzhu/copier`](https://github.com/jinzhu/copier) and a JSON round trip on your own machine:

```bash
# Library benchmarks
go test -run '^$' -bench . -benchmem

# Comparison benchmarks
cd benchmarks && go test -run '^$' -bench . -benchmem
```

The comparison benchmarks live in a separate module under [benchmarks](./benchmarks/), so the libraries they compare against are never dependencies of go-mapper.

## Best Practices

1. **Register early.** Register all mappings during application startup, before serving traffic.
2. **Keep mapping functions pure.** No side effects makes them easy to test and safe to call concurrently.
3. **Always check errors.** Validation in mapping functions only helps if callers look at the result.
4. **Clear the registry in tests.** `ClearRegistry()` keeps tests isolated from each other.

## Contributing

Contributions are welcome! See [CONTRIBUTING.md](CONTRIBUTING.md) for how to set up, test and submit changes, and [CHANGELOG.md](CHANGELOG.md) for what changed in each release.

If you find a bug or have an idea, please [open an issue](https://github.com/mdnishat27/go-mapper/issues/new/choose).

## License

[Apache License 2.0](LICENSE)

---

If go-mapper saves you some boilerplate, consider giving it a ⭐ on GitHub!