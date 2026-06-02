# macOS Setup

## Prerequisites

- macOS 12+ (Monterey or later)
- Go 1.25+ (only needed for building from source)

## Installation

### Option 1: Install script (recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/rkriad585/neovector/main/installer.sh | sh
```

Restart your shell, then:

```bash
neovector --help
```

### Option 2: Go install

```bash
go install github.com/rkriad585/neovector@latest
```

### Option 3: Download binary

Download from [releases](https://github.com/rkriad585/neovector/releases).

```bash
chmod +x neovector-darwin-arm64
sudo mv neovector-darwin-arm64 /usr/local/bin/neovector
```

## Build from source

```bash
git clone https://github.com/rkriad585/neovector.git
cd neovector
make build
```

## Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/rkriad585/neovector/main/installer.sh | sh -s -- --selfuninstall
```
