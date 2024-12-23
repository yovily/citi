# Go Module Template

A minimal, best-practices template for developing Go modules. This template provides a standardized starting point for creating reusable, high-quality Go modules that can be shared with the community.

[![Go Reference](https://pkg.go.dev/badge/github.com/username/modulename.svg)](https://pkg.go.dev/github.com/username/modulename)
[![Go Report Card](https://goreportcard.com/badge/github.com/username/modulename)](https://goreportcard.com/report/github.com/username/modulename)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

## Features

- Minimal, clean structure following Go best practices
- Built-in testing framework with examples
- Continuous Integration setup with GitHub Actions
- Comprehensive documentation structure
- Example implementations
- Development tools (linting, testing, building)

## Getting Started

### Using This Template

1. Click "Use this template" on GitHub
2. Clone your new repository
3. Update the module name:
   ```bash
   # Replace old module name
   find . -type f -name '*.go' -exec sed -i 's,github.com/username/modulename,github.com/your/module,g' {} +
   
   # Initialize module
   go mod init github.com/your/module
   
   # Install dependencies
   go mod tidy
   ```

### Prerequisites

- Go 1.21 or higher
- Make
- golangci-lint (optional, for development)

### Development Setup

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install dependencies
go mod download

# Run tests
make test

# Run linter
make lint
```

## Project Structure

```
.
├── .github/workflows/    # CI/CD configurations
├── configs/             # Configuration files
├── docs/               # Documentation
├── examples/           # Example implementations
├── test/              # Integration tests
├── module.go          # Core module code
├── options.go         # Configuration options
├── errors.go          # Error definitions
└── version.go         # Version information
```

## Best Practices for Go Module Development

### Code Organization
- Keep the public API minimal and focused
- Use interfaces for flexibility
- Follow standard Go project layout
- Keep dependencies minimal

### Documentation
- Write comprehensive godoc comments
- Include runnable examples
- Maintain a changelog
- Document all exported symbols

### Testing
- Write table-driven tests
- Include benchmarks
- Add examples as tests
- Test with race detector

### Version Management
- Use semantic versioning
- Tag releases
- Maintain compatibility
- Document breaking changes

### Error Handling
- Define custom error types
- Use error wrapping
- Return detailed errors
- Include error context

## Development Workflow

### Common Commands

```bash
# Run all tests
make test

# Run tests without integration tests
make test-short

# Run linter
make lint

# Generate coverage report
make coverage

# Run examples
make example-basic
make example-advanced

# Clean build artifacts
make clean
```

### Release Process

1. Update version in `version.go`
2. Update changelog
3. Run tests and linting:
   ```bash
   make test lint
   ```
4. Tag and push:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Module Best Practices Checklist

### Before Development
- [ ] Choose clear, descriptive module name
- [ ] Check for naming conflicts on pkg.go.dev
- [ ] Define clear module scope
- [ ] Plan public API carefully

### During Development
- [ ] Follow Go style guide
- [ ] Write comprehensive tests
- [ ] Add godoc examples
- [ ] Keep dependencies minimal
- [ ] Use Go modules
- [ ] Implement context support
- [ ] Add proper error handling
- [ ] Include benchmarks
- [ ] Document thoroughly

### Before Release
- [ ] Run all tests
- [ ] Check coverage
- [ ] Run linter
- [ ] Update documentation
- [ ] Tag version
- [ ] Update changelog
- [ ] Check backward compatibility

### After Release
- [ ] Monitor issues
- [ ] Respond to bug reports
- [ ] Track dependencies
- [ ] Plan future versions
- [ ] Maintain compatibility