# neovector

A fast CLI tool for seamless conversion between images and their numerical vector representations, rewritten in Go.

## Installation

```bash
go install github.com/rkriad585/neovector@latest
```

Or build from source:

```bash
git clone https://github.com/rkriad585/neovector.git
cd neovector
go build -o neovector.exe .
```

## Usage

### Convert an image to a vector

```bash
neovector convert to-vector <input> <output> [--format txt|json]
```

**Examples:**
```bash
neovector convert to-vector image.png vector.txt --format txt
neovector convert to-vector image.jpg vector.json --format json
```

### Convert a vector back to an image

```bash
neovector convert to-image <input> <output> <width> <height> [--format txt|json]
```

**Example:**
```bash
neovector convert to-image vector.txt restored.png 1920 1080
```

### Check dimensions

```bash
neovector check [--image <path>] [--vector <path>] [--format txt|json]
```

**Examples:**
```bash
neovector check --image image.png
neovector check --vector vector.txt
neovector check --image image.png --vector vector.txt
```

### Help

```bash
neovector --help
neovector convert --help
neovector convert to-vector --help
neovector convert to-image --help
neovector check --help
```

## License

MIT License. See [LICENSE](LICENSE) for details.
