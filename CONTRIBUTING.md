# Contributing to skube

Thank you for your interest in contributing to **skube**! We welcome contributions from the community to help make Kubernetes management more intuitive and accessible.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Reporting Bugs](#reporting-bugs)
- [Suggesting Features](#suggesting-features)

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment. We expect all contributors to:

- Be respectful and considerate
- Welcome newcomers and help them get started
- Focus on constructive feedback
- Accept responsibility and apologize when mistakes are made

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/skube.git
   cd skube
   ```
3. **Add the upstream repository**:
   ```bash
   git remote add upstream https://github.com/geminal/skube.git
   ```

## Development Setup

### Prerequisites

- **Go 1.21+** - [Install Go](https://golang.org/doc/install)
- **kubectl** - [Install kubectl](https://kubernetes.io/docs/tasks/tools/)
- **Access to a Kubernetes cluster** (for integration testing)

### Building from Source

```bash
# Install dependencies
go mod download

# Build the binary
go build -o skube cmd/skube/main.go

# Run tests
go test ./...

# Run with coverage
go test -cover ./...
```

## How to Contribute

### Types of Contributions

We welcome various types of contributions:

- **Bug fixes** - Help us squash bugs!
- **New features** - Add support for more Kubernetes resources or commands
- **Documentation** - Improve README, add examples, fix typos
- **Tests** - Increase test coverage
- **Performance improvements** - Make skube faster
- **Natural language improvements** - Add more synonyms or better parsing logic

### Workflow

1. **Create a branch** for your work:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our [coding standards](#coding-standards)

3. **Test your changes**:
   ```bash
   go test ./...
   ```

4. **Commit your changes** with clear, descriptive messages:
   ```bash
   git commit -m "feat: add support for StatefulSets"
   ```

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request** on GitHub

## Coding Standards

### Go Style

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` to format your code
- Run `go vet` to catch common mistakes
- Keep functions focused and under 50 lines when possible
- Use meaningful variable and function names

### Code Organization

- **Parser logic** goes in `internal/parser/`
- **Executor logic** goes in `internal/executor/`
- **Completion logic** goes in `internal/completion/`
- **Help/Version logic** goes in `internal/help/`
- **Configuration** goes in `internal/config/`
- **Tests** should be in the same package with `_test.go` suffix

### Adding New Commands

When adding a new command:

1. Add the command alias to `commandAliases` map in `parser.go`
2. Add parsing logic if needed in `parseCommand()`
3. Add executor handler in `executor.go`
4. Add test cases in `internal/parser/`
5. Update `COMMAND_REFERENCE.md` with examples

### Adding New Resources

When adding support for a new Kubernetes resource:

1. Add resource aliases to `resourceAliases` and `getCommandMap` in `parser.go`
2. Add handler function in `executor.go`
3. Add case in `ExecuteCommand()` switch statement
4. Add test cases
5. Update documentation

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/parser/...

# Run with verbose output
go test -v ./...
```

### Writing Tests

- Write table-driven tests when possible
- Test both success and error cases
- Use descriptive test names: `TestParseNaturalLanguage_WithStopWords`
- Aim for >70% code coverage

Example test structure:

```go
func TestYourFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {"description", "input", "expected"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := YourFunction(tt.input)
            if result != tt.expected {
                t.Errorf("got %v, want %v", result, tt.expected)
            }
        })
    }
}
```

## Pull Request Process

1. **Update documentation** if you're adding features or changing behavior
2. **Add tests** for new functionality
3. **Ensure all tests pass**: `go test ./...`
4. **Update COMMAND_REFERENCE.md** if adding new commands
5. **Write a clear PR description** explaining:
   - What problem does this solve?
   - What changes were made?
   - How was it tested?

### PR Title Format

Use conventional commit format:

- `feat: add support for CronJobs`
- `fix: correct namespace parsing for logs command`
- `docs: update README with installation instructions`
- `test: add tests for natural language parsing`
- `refactor: simplify parseCommand logic`
- `perf: optimize resource lookup`

### Review Process

- Maintainers will review your PR within a few days
- Address any feedback or requested changes
- Once approved, a maintainer will merge your PR

## Reporting Bugs

Found a bug? Please open an issue with:

- **Clear title** describing the problem
- **Steps to reproduce** the issue
- **Expected behavior** vs **actual behavior**
- **Environment details**: OS, Go version, kubectl version
- **Error messages** or logs if applicable

Example:

```markdown
**Bug**: `skube logs myapp` fails with namespace error

**Steps to Reproduce**:
1. Run `skube logs myapp in production`
2. Observe error

**Expected**: Should show logs from app
**Actual**: Error: "namespace not found"

**Environment**:
- OS: macOS 14.0
- Go: 1.21.3
- kubectl: 1.28.0
```

## Suggesting Features

Have an idea? Open an issue with:

- **Clear description** of the feature
- **Use case**: Why is this useful?
- **Proposed solution** (if you have one)
- **Alternatives considered**

## Questions?

- Open a [GitHub Discussion](https://github.com/geminal/skube/discussions)
- Check existing [Issues](https://github.com/geminal/skube/issues)

---

Thank you for contributing to skube! 🚀
