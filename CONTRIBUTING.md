# Contributing to SARC-NG

Thank you for your interest in contributing to SARC-NG! We welcome contributions from the community.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/sarc-ng.git
   cd sarc-ng
   ```
3. **Set up the development environment**:
   ```bash
   make setup
   ```

## Development Workflow

### Before You Start
- Check existing [issues](https://github.com/tecmx/sarc-ng/issues) and [pull requests](https://github.com/tecmx/sarc-ng/pulls)
- For major changes, please open an issue first to discuss the proposed changes

### Making Changes

1. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** following our coding standards:
   - Write clean, readable code
   - Follow Go conventions and best practices
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**:
   ```bash
   make test           # Run tests
   make lint           # Run linters
   make build          # Build the project
   ```

4. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

### Commit Message Format
Use conventional commits format:
- `feat:` for new features
- `fix:` for bug fixes
- `docs:` for documentation changes
- `test:` for test additions/changes
- `refactor:` for code refactoring
- `chore:` for maintenance tasks

### Submitting Changes

1. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

2. **Create a Pull Request**:
   - Go to the [SARC-NG repository](https://github.com/tecmx/sarc-ng)
   - Click "New Pull Request"
   - Select your branch and provide a clear description

## Code Standards

### Go Code
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Write comprehensive tests
- Add comments for complex logic
- Follow the project's architecture patterns

### API Documentation
- Update OpenAPI specification for API changes
- Include examples in API documentation
- Test API endpoints thoroughly

### Documentation
- Update relevant documentation for changes
- Use clear, concise language
- Include code examples where helpful

## Testing

- Write unit tests for new functions
- Include integration tests for API endpoints
- Ensure all tests pass before submitting PR
- Aim for good test coverage

## Questions or Issues?

- **Bug Reports**: [Create an issue](https://github.com/tecmx/sarc-ng/issues/new)
- **Feature Requests**: [Start a discussion](https://github.com/tecmx/sarc-ng/discussions)
- **Questions**: Check existing [discussions](https://github.com/tecmx/sarc-ng/discussions)

## Code of Conduct

Please be respectful and professional in all interactions. We're committed to providing a welcoming environment for all contributors.

## License

By contributing to SARC-NG, you agree that your contributions will be licensed under the same license as the project.

---

**Thank you for contributing to SARC-NG!** 🎉
