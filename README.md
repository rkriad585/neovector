# neovector

A fast, cross-platform CLI tool for converting between raster images and their
numerical vector representations (flat lists of RGB pixel values).

## Features

- **Image to vector** — convert any image to a numerical vector (txt, json, csv, bin)
- **Vector to image** — reconstruct images from vectors (with or without manual dimensions)
- **Metadata headers** — vector files embed width/height for automatic reconstruction
- **Batch processing** — glob patterns and recursive directory scanning
- **Grayscale & channel selection** — extract specific color channels
- **JPEG quality control** — set output quality for JPEG reconstruction
- **Dimension analysis** — inspect images and find possible dimensions
- **Self-update** — update to the latest version with `neovector self-update`
- **Self-uninstall** — fully remove neovector with `neovector self-uninstall`
- **13 Color Themes** — switch themes with `neovector config theme <name>`
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
# TXT format (with metadata header)
neovector convert to-vector image.png vector.txt

# JSON, CSV, or Binary
neovector convert to-vector image.jpg vector.json
neovector convert to-vector image.png vector.csv
neovector convert to-vector image.png vector.bin
```

### Convert vector to image

```bash
# With header (dimensions auto-detected)
neovector convert to-image vector.txt restored.png

# Without header (manual dimensions)
neovector convert to-image vector.txt restored.png 1920 1080

# JPEG with quality control
neovector convert to-image vector.json photo.jpg 1920 1080 --quality 85
```

### Batch conversion

```bash
# Convert all PNGs in batch
neovector convert to-vector "images/*.png" ./vectors --recursive

# Reconstruct all vectors with headers
neovector convert to-image "vectors/*.txt" ./images --recursive
```

### Grayscale and channels

```bash
# Single-channel grayscale
neovector convert to-vector photo.jpg gray.txt --grayscale

# Extract red channel only
neovector convert to-vector photo.jpg red.txt --channels r

# Extract red+green channels
neovector convert to-vector photo.jpg rg.txt --channels rg
```

### Check dimensions

```bash
neovector check --image image.png
neovector check --vector vector.txt
neovector check --image image.png --vector vector.txt
```

### Global flags

```bash
# Persistent flags available on all commands
neovector --output-dir "~/myoutput" convert to-vector image.png vector.txt
neovector --format csv convert to-vector image.png vector.csv
```

### Self-update

```bash
neovector self-update
neovector self-update --proxy http://proxy:8080
```

### Configuration

```bash
# View current config
neovector config

# Set values directly
neovector config --format csv
neovector config --output-dir "~/projects/vectors"
neovector config --proxy http://proxy:8080

# Interactive TUI editor
neovector config --edit

# Theme management
neovector config theme              # Show current theme + all themes
neovector config theme <name>       # Switch theme directly
neovector config theme --edit       # Interactive theme picker
```

### Self-uninstall

```bash
neovector self-uninstall
```

### Get version

```bash
neovector version
neovector --version
```

## Configuration

Config file: `~/.config/neostore/neovector/config.toml`

```toml
[general]
default_format = "txt"
output_dir = ""

[network]
proxy = ""

[theme]
name = "sunny_beach_day"
mode = "dark"
```

History log: `~/.config/neostore/neovector/history.log`

Output files default to: `~/Downloads/neostore/neovector/`
(override with `output_dir` in config or `--output-dir` flag)

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
│   ├── check.go            # dimension analysis
│   ├── config.go           # config + config theme subcommand
│   ├── completion.go       # Shell completion generation
│   ├── selfupdate.go       # self-update command
│   └── selfuninstall.go    # self-uninstall subcommand
├── internal/
│   ├── banner/             # Startup banner
│   ├── config/             # Config dir, config.toml, output dir
│   ├── log/                # history.log writer
│   ├── theme/              # Color theme system (13 themes)
│   ├── uninstall/          # Self-uninstall logic
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
