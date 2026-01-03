# Agent Instructions

Instructions for AI coding agents when working on this project.

## Agent Organization

When creating agent instruction files:

- The main file should always be named `AGENTS.md`
- Create a `CLAUDE.md` symlink pointing to `AGENTS.md` for compatibility with Claude Code

## Project Overview

PublicSuffix for Go is a Go domain name parser based on the Public Suffix List. It provides flexible parsing, validation, and domain extraction with support for private domains, IDNA/Punycode, and multiple list configurations.

## Key Documentation

- **[README.md](README.md)** - Library overview, features, usage examples, and API reference
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Contribution guidelines, commit format, testing approach
- **[Public Suffix List](https://publicsuffix.org/)** - Official Public Suffix List documentation
- **[Go Documentation](https://pkg.go.dev/github.com/weppos/publicsuffix-go/publicsuffix)** - API documentation

## Project-Specific Context

### Code Style Notes

- Use Conventional Commits format (see [CONTRIBUTING.md](CONTRIBUTING.md#commit-message-guidelines))
- Follow Go best practices and idiomatic Go patterns
- Do not include AI attribution in commit messages or code comments
- Tests are mandatory for all changes

## Project Structure

```
publicsuffix/       # Main package
net/publicsuffix/   # Compatibility adapter for golang.org/x/net/publicsuffix
cmd/                # CLI tools
fixtures/           # Test data
```
