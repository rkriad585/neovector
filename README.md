# neovector

A fast, cross-platform CLI tool for converting between raster images and their
numerical vector representations (flat lists of RGB pixel values).

## Features

- **Image to vector** — convert any image to a numerical vector (txt/json)
- **Vector to image** — reconstruct images from vectors
- **Dimension analysis** — inspect images and find possible dimensions
- **Cross-platform** — Windows, Linux, macOS binaries
- **Portable** — single static binary, no runtime dependencies
- **Configurable** — config.toml at `~/.config/neostore/neovector/`
- **Operation log** — every conversion logged with timestamp

## Installation

### From source

```bash
git clone https://github.com/rkriad585/neovector.git
cd neovector
go build -o neovector .
```

### Via install script

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/rkriad585/neovector/main/installer.ps1 | iex
```

**Linux/macOS:**
```bash
curl -fsSL https://raw.githubusercontent.com/rkriad585/neovector/main/installer.sh | sh
```

### With Go install

```bash
go install github.com/rkriad585/neovector@latest
```

### Docker

```bash
docker build -t rkriad585/neovector .
docker run --rm rkriad585/neovector --help
```

## Usage

### Convert image to vector

```bash
neovector convert to-vector image.png vector.txt
neovector convert to-vector image.jpg vector.json --format json
```

### Convert vector to image

```bash
neovector convert to-image vector.txt restored.png 1920 1080
```

### Check dimensions

```bash
neovector check --image image.png
neovector check --vector vector.txt
neovector check --image image.png --vector vector.txt
```

### Get version

```bash
neovector --version
```

## Configuration

Config file: `~/.config/neostore/neovector/config.toml`

```toml
[general]
default_format = "txt"
```

History log: `~/.config/neostore/neovector/history.log`

Output files default to: `~/Downloads/neostore/neovector/`

## Build

```bash
# Build for current platform
make build

# Cross-compile all platforms
make release

# Build with Docker
make docker-build

# Run tests
make test

# Run linter
make lint
```

### Manual build

```bash
go build -trimpath -ldflags "-s -w" -o neovector .
```

## Project Structure

```
neovector/
├── main.go                 # Entry point, ldflags injection
├── .version                # Current version
├── cmd/
│   ├── root.go             # Root cobra command, banner, config init
│   ├── convert.go          # to-vector, to-image subcommands
│   └── check.go            # dimension analysis
├── internal/
│   ├── banner/             # Startup banner
│   ├── config/             # Config dir, config.toml, output dir
│   ├── log/                # history.log writer
│   └── version/            # Build metadata vars
├── build.ps1               # Windows cross-compile script
├── build.sh                # Linux/macOS cross-compile script
├── installer.ps1           # Windows installer
├── installer.sh            # Linux/macOS installer
├── Dockerfile              # Multi-stage Docker build
├── docker-compose.yml      # Docker Compose
├── Makefile                # Build automation
└── CMakeLists.txt          # CMake build
```

## License

MIT License. See [LICENSE](LICENSE).
