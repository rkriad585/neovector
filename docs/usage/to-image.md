# to-image — Convert Vector to Image

Reconstruct an image from a numerical vector file. Width and height
can be omitted if the vector has a metadata header (txt and csv
formats include headers by default).

## Usage

```bash
neovector convert to-image <input> <output> [width] [height] [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `input`  | Path to the vector file or glob pattern |
| `output` | Output image path (or directory for batch mode) |
| `width`  | Image width in pixels (optional if header present) |
| `height` | Image height in pixels (optional if header present) |

## Flags

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `--format` | | `txt` | Input vector format: `txt`, `json`, `csv`, `bin` |
| `--quality` | | `95` | JPEG output quality (1-100) |
| `--overwrite` | `-o` | `false` | Overwrite output file if it exists |
| `--verbose` | `-V` | `false` | Show detailed conversion progress |
| `--recursive` | `-r` | `false` | Process directories recursively (batch mode) |

## Examples

### Reconstruct from txt (with header)

```bash
neovector convert to-image vector.txt restored.png
```
Width and height are read from the vector's metadata header.

### Reconstruct from txt (without header)

```bash
neovector convert to-image vector.txt restored.png 1920 1080
```

### Reconstruct from JSON

```bash
neovector convert to-image vector.json restored.png 800 600
```

### Reconstruct from binary

```bash
neovector convert to-image vector.bin restored.png 640 480
```

### JPEG with custom quality

```bash
neovector convert to-image vector.json photo.jpg 1920 1080 --quality 80
```

### Batch convert all vectors in a directory

```bash
neovector convert to-image "vectors/*.txt" ./images --recursive
```

### Find dimensions first with check

```bash
neovector check --vector vector.txt
```

The check command shows all possible (width, height) factor pairs,
and also displays header dimensions if present.

## Notes

- The vector size must exactly equal `width * height * 3` (unless selected channels)
- When width/height are provided explicitly, they override the header
- Output format is determined by file extension (`.jpg`/`.jpeg` → JPEG, else PNG)
- In batch mode, vectors must have headers (or use explicit width/height flags)
