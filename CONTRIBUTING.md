# Contributing to SARC-NG

Thank you for your interest in contributing to SARC-NG! We appreciate your help in making this project better.

## Getting Started

```bash
# Fork and clone
git clone https://github.com/YOUR-USERNAME/sarc-ng.git
cd sarc-ng
make setup
```

## Development Workflow

1. **Create a branch**
   ```bash
   git checkout -b feature/your-feature
   ```

2. **Make changes** - Follow Go conventions and write tests

3. **Run checks**
   ```bash
   make test
   make lint
   make build
   ```

4. **Commit**
   ```bash
   git commit -m "feat: add your feature"
   ```

5. **Push and create PR**
   ```bash
   git push origin feature/your-feature
   ```

## Commit Format

Use conventional commits:
- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation
- `test:` - Tests
- `refactor:` - Code refactoring
- `chore:` - Maintenance

## Code Standards

- Follow Go formatting (`gofmt`)
- Write meaningful names
- Add tests for new code
- Update documentation
- Comment complex logic

## Pull Request Process

1. **Update documentation** for any new features
2. **Add tests** for new functionality (aim for >80% coverage)
3. **Update CHANGELOG** if applicable
4. **Link related issues** in the PR description
5. **Request review** from at least one maintainer
6. **Address feedback** promptly and respectfully

### PR Template

When creating a PR, please include:
- Description of changes
- Related issue(s)
- Type of change (feature, bugfix, docs, etc.)
- Testing performed
- Screenshots (if UI changes)

## Issue Reporting

### Bug Reports

When reporting bugs, include:
- Go version
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Error messages and logs
- Minimal reproducible example

### Feature Requests

For feature requests, describe:
- The problem you're trying to solve
- Proposed solution
- Alternative solutions considered
- Impact on existing functionality

## Code of Conduct

### Our Pledge

We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity and expression, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards

**Positive behavior includes:**
- Using welcoming and inclusive language
- Being respectful of differing viewpoints
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards others

**Unacceptable behavior includes:**
- Harassment, trolling, or derogatory comments
- Personal or political attacks
- Public or private harassment
- Publishing others' private information
- Other conduct which could reasonably be considered inappropriate

### Enforcement

Instances of abusive, harassing, or otherwise unacceptable behavior may be reported by contacting the project team. All complaints will be reviewed and investigated promptly and fairly.

## Security Vulnerabilities

**Do NOT create public issues for security vulnerabilities.**

If you discover a security vulnerability, please follow these steps:

1. **Email the security team** at security@example.com (or create a private security advisory on GitHub)
2. **Include details:**
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)
3. **Wait for response** - We aim to respond within 48 hours
4. **Coordinate disclosure** - We'll work with you on a coordinated disclosure timeline

We appreciate your responsible disclosure and will credit you in the security advisory (unless you prefer to remain anonymous).

## Development Setup

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Make
- Git

### Environment Setup

```bash
# Install dependencies
go mod download

# Start development environment
make docker-up

# Run tests
make test

# Run linters
make lint
```

### Testing

```bash
# Run all tests
make test

# Run specific test
go test -v ./internal/domain/building/...

# Run with coverage
make coverage

# Run integration tests
make test-integration
```

## Documentation

- Update README.md for user-facing changes
- Add inline code comments for complex logic
- Update API documentation (OpenAPI spec)
- Add examples for new features

## Questions?

- **Bugs**: [Create an issue](https://github.com/tecmx/sarc-ng/issues)
- **Features**: [Start a discussion](https://github.com/tecmx/sarc-ng/discussions)
- **Security**: Email security@example.com
- **General**: [Community discussions](https://github.com/tecmx/sarc-ng/discussions)

## License

By contributing to SARC-NG, you agree that your contributions will be licensed under the same license as the project.

---

Thank you for contributing! 🎉
