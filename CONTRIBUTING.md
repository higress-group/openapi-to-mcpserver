# Contributing Guide

Thank you for your interest in the **openapi-to-mcpserver** project! This project is a tool that converts OpenAPI specifications to MCP (Model Context Protocol) server configurations.

We welcome various forms of contributions, including but not limited to:

- 🐛 Report and fix bugs
- ✨ Add new features
- 📚 Improve documentation
- 🧪 Add or improve tests
- 🎨 Code optimization and refactoring

## 📋 Table of Contents

- Quick Start
- Development Environment Setup
- Project Structure
- Code Standards
- Testing Guidelines
- Submitting Pull Requests
- Issue Reporting
- License

## 🚀 Quick Start

### Environment Requirements

- Go 1.23.0 or higher
- Git

### Clone and Setup Project

```bash
# Clone the project
git clone https://github.com/higress-group/openapi-to-mcpserver.git
cd openapi-to-mcpserver

# Install dependencies
go mod download

# Verify installation
go build ./cmd/openapi-to-mcp
```

### Run Example

```bash
# Convert using example OpenAPI specification
go run ./cmd/openapi-to-mcp --input test/petstore.json --output example-mcp.yaml --server-name petstore

# View help information
go run ./cmd/openapi-to-mcp --help
```

## 🛠️ Development Environment Setup

### 1. Install Go

Ensure you have Go 1.23.0 or higher installed:

```bash
# Check Go version
go version

# If you need to upgrade Go, you can use:
# macOS (using Homebrew)
brew install go

# Linux
# Download and install latest version from https://golang.org/dl/
```

### 2. Setup Development Environment

```bash
# Clone the project
git clone https://github.com/higress-group/openapi-to-mcpserver.git
cd openapi-to-mcpserver

# Install dependencies
go mod tidy

# Verify build
go build ./...

# Run tests
go test ./...
```

### 3. IDE Configuration

Recommended IDEs and plugins:

- **VS Code** + Go extension
- **GoLand** (JetBrains)
- **Vim/Neovim** + vim-go plugin

### 4. Common Development Commands

```bash
# Format code
go fmt ./...

# Run lint checks
go vet ./...

# Run tests
go test ./...

# Run tests and generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Build binary
go build ./cmd/openapi-to-mcp

# Install locally
go install ./cmd/openapi-to-mcp
```

## 🏗️ Project Structure

```text
openapi-to-mcpserver/
├── cmd/openapi-to-mcp/          # Main program entry point
│   └── main.go                 # CLI application
├── pkg/                        # Core packages
│   ├── converter/              # OpenAPI to MCP conversion logic
│   │   ├── converter.go        # Main converter implementation
│   │   └── converter_test.go   # Converter tests
│   ├── models/                 # Data model definitions
│   │   └── mcp_config.go       # MCP configuration structures
│   └── parser/                 # OpenAPI parser
│       └── parser.go           # OpenAPI document parsing
├── test/                       # Test data and expected outputs
│   ├── *.json                  # OpenAPI test specifications
│   ├── expected-*.yaml         # Expected MCP configuration outputs
│   └── template.yaml           # Test template
├── go.mod                      # Go module definition
├── go.sum                      # Go dependency checksums
├── README.md                   # Project documentation
├── LICENSE                     # Apache 2.0 license
└── CONTRIBUTING.md             # This contributing guide
```

## 📝 Code Standards

### Go Code Standards

We follow Go's official code standards and best practices:

#### 1. Code Formatting

Use `go fmt` for code formatting:

```bash
go fmt ./...
```

#### 2. Code Quality Checks

```bash
# Run go vet to check for potential issues
go vet ./...

# Use golangci-lint (if installed)
golangci-lint run
```

#### 3. Naming Conventions

- **Package names**: Use lowercase, short names
- **Functions/Methods**: Use camelCase, capitalize first letter for public
- **Variables**: Use camelCase, lowercase first letter for private
- **Constants**: Use ALL_CAPS with underscores

#### 4. Comment Standards

```go
// Package parser provides functionality to parse OpenAPI specifications.
package parser

// ParseFile parses an OpenAPI specification from a file.
func ParseFile(filename string) error {
    // TODO: Add file validation
    return nil
}
```

### Commit Message Standards

Commit messages should be clear and concise, following this format:

```text
<type>(<scope>): <subject>

<body>

<footer>
```

#### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation update
- `style`: Code style adjustment
- `refactor`: Code refactoring
- `test`: Test-related changes
- `chore`: Build process or tooling configuration updates

#### Example

```text
feat: add support for OpenAPI 3.1

Add support for parsing OpenAPI 3.1 specifications with new features
like webhooks and improved schema validation.

Closes #123
```

## 🧪 Testing Guidelines

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests for specific package
go test ./pkg/converter

# Run tests with detailed output
go test -v ./...

# Run tests and generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Tests

1. **Unit Tests**: Write unit tests for each function
2. **Integration Tests**: Test interactions between components
3. **End-to-End Tests**: Use `TestEndToEndConversion` pattern

#### Test File Structure

```go
package converter

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestConvertOperation(t *testing.T) {
    // Arrange
    p := parser.NewParser()
    c := NewConverter(p, models.ConvertOptions{})

    // Act
    result, err := c.convertOperation("/pets", "get", operation)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "listPets", result.Name)
}
```

### Test Coverage Goals

- Target test coverage: > 80%
- Core functionality coverage: > 90%

## 📤 Submitting Pull Requests

### 1. Preparation

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Or fix branch
git checkout -b fix/issue-number
```

### 2. Development Process

1. Write code following coding standards
2. Add or update tests
3. Run tests to ensure they pass
4. Update documentation (if needed)
5. Commit changes

### 3. Submit PR

1. **Push branch** to your fork
2. **Create PR** and include:
   - Clear title
   - Detailed description
   - Related issue references
   - Test screenshots (if applicable)

### 4. PR Template

Please use the following template:

```markdown
## Description

Brief description of your changes.

## Type

- [ ] 🐛 Bug fix
- [ ] ✨ New feature
- [ ] 📚 Documentation update
- [ ] 🎨 Code refactoring
- [ ] 🧪 Test updates
- [ ] 🔧 Build/tooling updates

## Checklist

- [ ] My code follows the project's coding standards
- [ ] I have added necessary tests
- [ ] All tests are passing
- [ ] I have updated relevant documentation
- [ ] My changes do not break existing functionality

## Related Issues

Closes #123
```

## 🐛 Issue Reporting

### Bug Reports

If you find a bug, please:

1. **Check existing issues** to see if it has already been reported
2. **Create a new issue** and provide:
   - Clear title
   - Detailed description
   - Steps to reproduce
   - Expected behavior
   - Actual behavior
   - Environment information

### Feature Requests

For new feature requests, please:

1. **Describe the problem** that this feature will solve
2. **Describe the solution** you would like to see implemented
3. **Consider alternatives** to the proposed solution

## 📄 License

By contributing code, you agree that your contributions will be licensed under the [Apache License 2.0](LICENSE).

## 🙋‍♂️ Getting Help

If you have any questions:

1. Check [README.md](README.md) for basic information
2. Review existing issues and discussions
3. Create a new issue or discussion

Thank you for your contributions! 🎉
