# Building from Source

## Prerequisites

- Go 1.25+
- Git

## Quick Build

```bash
git clone https://github.com/rkriad585/neovector.git
cd neovector
go build -o neovector .
```

## Makefile Targets

| Target | Description |
|--------|-------------|
| `make build` | Build for current platform |
| `make run` | Build and run with `ARGS` |
| `make test` | Run all tests |
| `make lint` | Run go vet |
| `make format` | Format code with go fmt |
| `make clean` | Remove build artifacts |
| `make release` | Cross-compile all 6 platforms |
| `make install` | Install to `$GOPATH/bin` |
| `make docker-build` | Build Docker image |
| `make docker-run` | Run Docker container |

## Cross-Compiling

Using the build scripts:

```bash
# Windows
.\build.ps1

# Linux/macOS
chmod +x build.sh
./build.sh
```

Output goes to `bin/`:

```
bin/neovector-windows-amd64.exe
bin/neovector-windows-arm64.exe
bin/neovector-linux-amd64
bin/neovector-linux-arm64
bin/neovector-darwin-amd64
bin/neovector-darwin-arm64
```

## Docker

```bash
docker build -t rkriad585/neovector .
docker run --rm rkriad585/neovector --help

# Mount volumes for file access
docker run --rm \
  -v $(pwd)/images:/data:ro \
  -v $(pwd)/output:/output \
  rkriad585/neovector convert to-vector /data/photo.jpg /output/vector.txt
```

## Testing

```bash
go test ./... -v

# With coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Versioning

The version is read from `.version` at build time. To create a release:

```bash
echo "v1.0.0" > .version
git add .version
git commit -m "chore: bump to v1.0.0"
git tag v1.0.0
git push --tags
```

The GitHub Actions workflow will build and publish all binaries automatically.
