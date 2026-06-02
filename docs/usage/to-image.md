# to-image — Convert Vector to Image

Reconstruct an image from a numerical vector file.

## Usage

```bash
neovector convert to-image <input> <output> <width> <height> [--format txt|json]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `input`  | Path to the vector file |
| `output` | Path for the output image |
| `width`  | Image width in pixels |
| `height` | Image height in pixels |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `txt` | Input vector format: `txt` or `json` |

## Examples

### Reconstruct from txt

```bash
neovector convert to-image vector.txt restored.png 1920 1080
```

### Reconstruct from JSON

```bash
neovector convert to-image vector.json restored.png 800 600 --format json
```

### Find dimensions first with check

```bash
neovector check --vector vector.txt
```

The check command shows all possible (width, height) factor pairs.

## Notes

- The vector size must exactly equal `width * height * 3`
- Output format is determined by file extension (`.jpg`/`.jpeg` → JPEG, else PNG)
- JPEG quality is fixed at 95
