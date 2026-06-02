# Contributing

Thank you for considering contributing to neovector.

## How to Contribute

1. Fork the repository
2. Create a feature branch: `git checkout -b feat/my-feature`
3. Commit your changes with conventional commits
4. Push to your fork and open a pull request

## Commit Convention

Use [conventional commits](https://www.conventionalcommits.org/):

```
feat: add new feature
fix: correct a bug
docs: update documentation
perf: improve performance
refactor: restructure code
test: add or fix tests
chore: maintenance tasks
```

## Development Setup

```bash
git clone https://github.com/rkriad585/neovector.git
cd neovector
go mod download
go build -o neovector .
```

## Code Style

- Run `go vet ./...` before committing
- Run `go fmt ./...` to format code
- Ensure all tests pass: `go test ./...`
- Keep the banner width at 56 characters if modifying it

## Pull Request Process

1. Update documentation if needed
2. Ensure the build passes
3. Link any related issues
4. Request review from maintainers

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
