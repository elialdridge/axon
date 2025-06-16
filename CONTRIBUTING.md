# Contributing to Axon

Thank you for your interest in contributing to Axon! This document provides guidelines and information for contributors.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Please read it before contributing.

## How to Contribute

### Reporting Bugs

1. **Check existing issues** first to avoid duplicates
2. **Use the bug report template** when creating new issues
3. **Provide detailed information**:
   - Operating system and version
   - Go version
   - Terminal emulator
   - Steps to reproduce
   - Expected vs actual behavior
   - Error messages or logs

### Suggesting Features

1. **Check existing feature requests** to avoid duplicates
2. **Use the feature request template**
3. **Explain the use case** and how it aligns with Axon's philosophy
4. **Consider the impact** on minimalism and modularity

### Contributing Code

#### Prerequisites

- Go 1.23 or later
- Git
- Familiarity with terminal applications
- Understanding of Axon's architecture

#### Setup Development Environment

```bash
# Fork the repository on GitHub
# Clone your fork
git clone https://github.com/YOUR_USERNAME/axon.git
cd axon

# Add upstream remote
git remote add upstream https://github.com/ORIGINAL_OWNER/axon.git

# Install dependencies
go mod download

# Install development tools
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Build and test
make build
make test
```

#### Development Workflow

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards

3. **Write tests** for new functionality

4. **Run tests and linting**:
   ```bash
   make test
   make lint
   ```

5. **Commit your changes** using conventional commits:
   ```bash
   git commit -m "feat: add new world generation feature"
   git commit -m "fix: resolve save file corruption issue"
   git commit -m "docs: update API documentation"
   ```

6. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Create a Pull Request** on GitHub

#### Coding Standards

##### Go Style Guidelines

- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use `gofmt` for formatting
- Use `goimports` for import organization
- Follow Go naming conventions
- Write self-documenting code

##### Documentation

- Document all public functions and types
- Use meaningful variable and function names
- Include examples in godoc comments
- Keep comments up to date with code changes

##### Error Handling

- Always handle errors explicitly
- Provide context in error messages
- Use wrapped errors where appropriate:
  ```go
  if err != nil {
      return fmt.Errorf("failed to save game: %w", err)
  }
  ```

##### Testing

- Write tests for all new functionality
- Maintain or improve test coverage
- Use table-driven tests where appropriate
- Mock external dependencies
- Test error conditions

##### Architecture Principles

1. **Extreme Minimalism**
   - Keep features simple and focused
   - Avoid feature creep
   - Question every addition

2. **Maximum Modularity**
   - Separate concerns clearly
   - Use interfaces for abstraction
   - Make components replaceable

3. **Terminal Native**
   - Design for terminal environments
   - Ensure cross-platform compatibility
   - Respect terminal limitations

#### Pull Request Guidelines

##### Before Submitting

- [ ] Code follows style guidelines
- [ ] All tests pass
- [ ] Documentation is updated
- [ ] Commit messages follow conventional format
- [ ] No breaking changes (or properly documented)
- [ ] Performance impact considered
- [ ] Security implications reviewed

##### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] Manual testing performed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added and passing
```

### Code Review Process

1. **Automated checks** must pass (CI/CD)
2. **Maintainer review** for code quality and architecture
3. **Community feedback** on significant changes
4. **Final approval** from project maintainers

#### Review Criteria

- Code quality and readability
- Adherence to architecture principles
- Test coverage and quality
- Documentation completeness
- Performance impact
- Security considerations
- Backward compatibility

## Development Guidelines

### Project Structure

```
axon/
├── main.go                 # Entry point
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── ai/                # AI client implementation
│   ├── game/              # Core game logic
│   ├── storage/           # Save/load functionality
│   └── ui/                # User interface components
├── pkg/                   # Public API (if any)
├── docs/                  # Additional documentation
├── scripts/               # Build and utility scripts
└── testdata/              # Test fixtures
```

### Adding New Features

#### 1. AI Integration Features

- Implement in `internal/ai/`
- Add model selection logic
- Include error handling and fallbacks
- Test with mock AI responses

#### 2. Game Mechanics

- Implement in `internal/game/`
- Maintain state consistency
- Add to save/load system
- Include comprehensive tests

#### 3. UI Components

- Implement in `internal/ui/`
- Follow monochrome design
- Ensure terminal compatibility
- Test across different screen sizes

#### 4. Storage Features

- Implement in `internal/storage/`
- Maintain backward compatibility
- Include migration logic
- Test data integrity

### Testing Guidelines

#### Unit Tests

```go
func TestGameEngine_ProcessAction(t *testing.T) {
    tests := []struct {
        name     string
        action   string
        state    *GameState
        expected string
        wantErr  bool
    }{
        {
            name:   "valid action",
            action: "look around",
            state:  newTestState(),
            expected: "test response",
            wantErr: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

#### Integration Tests

```go
// +build integration

func TestFullGameSession(t *testing.T) {
    // Test complete game workflow
}
```

#### Benchmark Tests

```go
func BenchmarkAIGeneration(b *testing.B) {
    for i := 0; i < b.N; i++ {
        // Benchmark AI generation
    }
}
```

### Documentation

#### Code Documentation

```go
// ProcessAction handles player actions and generates AI responses.
// It validates the action, updates game state, and returns the narrative response.
//
// Example:
//   response, err := engine.ProcessAction(state, "examine the door")
//   if err != nil {
//       log.Fatal(err)
//   }
func (e *Engine) ProcessAction(state *GameState, action string) (string, error) {
    // Implementation
}
```

#### README Updates

Update README.md when adding:
- New features
- Configuration options
- Installation methods
- Usage examples

## Release Process

### Version Numbering

We follow [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH`
- `MAJOR`: Breaking changes
- `MINOR`: New features (backward compatible)
- `PATCH`: Bug fixes (backward compatible)

### Release Checklist

- [ ] All tests passing
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version tagged
- [ ] Binaries built for all platforms
- [ ] Release notes prepared

## Getting Help

### Communication Channels

- **GitHub Discussions**: General questions and ideas
- **GitHub Issues**: Bug reports and feature requests
- **Discord**: Real-time community chat
- **Email**: Direct contact with maintainers

### Resources

- [Go Documentation](https://golang.org/doc/)
- [Bubble Tea Framework](https://github.com/charmbracelet/bubbletea)
- [Project Architecture](DOCUMENTATION.md#architecture)
- [API Documentation](DOCUMENTATION.md#api-documentation)

## Recognition

Contributors are recognized in:
- CONTRIBUTORS.md file
- Release notes
- GitHub contributor graph
- Special mentions for significant contributions

## License

By contributing to Axon, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Axon! Your efforts help make AI-driven gaming more accessible and enjoyable for everyone.

