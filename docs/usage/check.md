# check — Inspect Image and Vector Dimensions

Inspect image dimensions or find possible reconstruction dimensions
from a vector file.

## Usage

```bash
neovector check [--image <path>] [--vector <path>] [--format txt|json]
```

## Flags

| Flag | Alias | Description |
|------|-------|-------------|
| `--image` | `-i` | Path to an image file |
| `--vector` | `-v` | Path to a vector file |
| `--format` | | Vector file format: `txt` or `json` (default: `txt`) |

## Examples

### Check an image

```bash
neovector check --image photo.jpg
```

### Check a vector

```bash
neovector check --vector vector.txt
```

### Check both at once

```bash
neovector check --image photo.jpg --vector vector.txt
```

### Use with JSON format

```bash
neovector check --vector vector.json --format json
```

## Output

### Image check shows

- Dimensions (width x height)
- Total pixel count
- Expected RGB vector size

### Vector check shows

- Vector format and size
- All possible (width, height) pairs if the size is divisible by 3
