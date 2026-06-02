# Development Workflow

## Branch Strategy

- `main` — stable, release-ready
- `feat/*` — new features
- `fix/*` — bug fixes
- `docs/*` — documentation updates
- `chore/*` — maintenance tasks

## Conventional Commits

```
feat(cmd): add resize flag to to-image
fix(convert): handle empty vector file gracefully
docs(readme): add Docker usage examples
perf(core): pre-allocate pixel slice
```

## Code Review Checklist

- [ ] `go vet ./...` passes
- [ ] `go test ./...` passes
- [ ] No commented-out code
- [ ] Error messages are user-friendly
- [ ] New functionality has tests
- [ ] Banner width (56) preserved if modified
- [ ] Documentation updated

## Releasing

1. Update `.version` with new version
2. Commit: `chore: bump to v0.2.0`
3. Tag: `git tag v0.2.0`
4. Push: `git push --tags`
5. GitHub Actions builds and publishes release
