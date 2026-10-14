# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Benchmark results comparing go-mapper with a manual copy, reflection, `jinzhu/copier` and a JSON round trip in the README
- Contributing guide, security policy, changelog, and issue and pull request templates
- Linting with golangci-lint, test coverage reporting and Dependabot in CI

### Changed

- Rewrote the README

## [0.1.0] - 2026-10-14

### Added

- `RegisterMapper`, `Map` and `MapList` for type-safe mapping between any two types using generics
- `HasMapping`, `UnregisterMapper` and `ClearRegistry` for managing registered mappings
- `ErrMapperNotFound`, wrapped by `Map` and `MapList` when no mapping is registered
- Support for value, pointer and interface types, each registered as its own type pair
- A registry that is safe for concurrent use

[Unreleased]: https://github.com/mdnishat27/go-mapper/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/mdnishat27/go-mapper/releases/tag/v0.1.0