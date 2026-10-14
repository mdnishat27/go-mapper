# Contributing to go-mapper

Thanks for your interest in improving go-mapper! Bug reports, ideas, documentation fixes and code are all welcome.

## Reporting bugs and suggesting features

Please [open an issue](https://github.com/mdnishat27/go-mapper/issues/new/choose) and pick the matching template. For bugs, a small code sample that reproduces the problem helps the most.

For larger changes, especially anything that changes the public API, please open an issue to discuss the idea before writing the code.

Security problems should not be reported in public issues. See [SECURITY.md](SECURITY.md) instead.

## Development setup

You need Go 1.21 or later.

```bash
git clone https://github.com/<your-username>/go-mapper.git
cd go-mapper
go test ./...
```

## Before opening a pull request

Run these from the repository root and make sure they are clean:

```bash
gofmt -l .           # must print nothing
go vet ./...
go test -race ./...
golangci-lint run    # see https://golangci-lint.run/welcome/install/
```

`-race` needs cgo. On Windows that means a 64-bit gcc on your `PATH` (for example from [MSYS2](https://www.msys2.org/)). If you cannot run it locally, CI runs it on Linux for every pull request.

If you touch the benchmarks, also run them. They live in their own module so their dependencies never reach go-mapper's users:

```bash
cd benchmarks
go test -run '^$' -bench . -benchmem
```

## Pull request guidelines

- Keep each pull request focused on one change.
- Add or update tests for any behavior change. New public functions also need a doc comment, and ideally a runnable example in `example_test.go`.
- Update the README if you change how something is used.
- Add a line under `[Unreleased]` in [CHANGELOG.md](CHANGELOG.md).
- go-mapper has no dependencies, and new ones are unlikely to be accepted.

## License

By contributing, you agree that your contributions are licensed under the [Apache License 2.0](LICENSE).