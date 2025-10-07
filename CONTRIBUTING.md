# Contributing to SARC-NG

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

## Questions?

- **Bugs**: [Create an issue](https://github.com/tecmx/sarc-ng/issues)
- **Features**: [Start a discussion](https://github.com/tecmx/sarc-ng/discussions)

Thank you for contributing! 🎉
