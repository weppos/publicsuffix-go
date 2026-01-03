# Contributing to PublicSuffix for Go

Thank you for your interest in contributing to PublicSuffix for Go!

## Development Workflow

1. Fork and clone the repository
2. Create a branch: `git checkout -b feature/your-feature`
3. Make your changes
4. Run tests: `go test ./...`
5. Commit using Conventional Commits format (see below)
6. Push and create a pull request

## Commit Message Guidelines

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification for commit messages.

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

Use lowercase for `<type>`. Capitalize the first letter of `<subject>` (sentence-style). See examples below.

### Type

- **feat**: A new feature
- **fix**: A bug fix
- **chore**: Other changes that don't modify src or test files
- **docs**: Documentation only changes
- **style**: Code style changes (formatting, etc.)
- **refactor**: Code change that neither fixes a bug nor adds a feature
- **perf**: Performance improvement
- **test**: Adding or updating tests
- **build**: Changes to build system or external dependencies
- **ci**: Changes to CI configuration files and scripts

### Scope

- **parser**: List parsing and rule handling
- **domain**: Domain parsing and validation
- **list**: List management and operations
- **adapter**: golang.org/x/net/publicsuffix adapter
- **cmd**: Command-line tools

### Examples

```bash
feat(parser): add support for wildcard exceptions

fix(domain): handle empty labels correctly

docs: update usage examples in README

refactor(list): simplify rule lookup logic
```

### Breaking Changes

Add `BREAKING CHANGE:` in the footer:

```
feat(domain): change Parse return signature

BREAKING CHANGE: Parse now returns error instead of panic.
Update code to handle returned errors.
```

## Testing

### Running Tests

```bash
go test ./...                    # Run all tests
go test ./publicsuffix           # Run package tests
go test -v ./...                 # Verbose output
go test -race ./...              # Run with race detector
```

### Writing Tests

- Unit tests should be in `*_test.go` files
- Use table-driven tests for multiple cases
- Test both happy paths and error cases
- Include tests for edge cases (empty strings, unicode, etc.)

## Pull Request Process

1. Update documentation for API changes
2. Add tests for new features or bug fixes
3. Ensure all tests pass: `go test ./...`
4. Use Conventional Commits format
5. Provide clear PR description with context
6. Reference related issues if applicable

## Code Style

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` to format code
- Run `go vet` to catch common issues
- Add godoc comments for exported functions/types
- Keep functions focused and small

### Go-Specific Guidelines

- Use idiomatic Go patterns
- Prefer composition over inheritance
- Handle errors explicitly, don't ignore them
- Use meaningful variable names
- Avoid global state when possible
- Use interfaces for abstraction

### Documentation

- Add godoc comments for all exported symbols
- Include usage examples in documentation
- Update README.md for user-facing changes
- Keep inline comments minimal - prefer self-documenting code

## Questions?

Open an issue for questions, feature discussions, or bug reports.

Thank you for contributing!
