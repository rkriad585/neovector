# to-vector — Convert Image to Vector

Convert a raster image into a flat list of RGB pixel values.
The output file includes a metadata header with width, height, and channel
count so that `to-image` can reconstruct the image without manual dimensions.

## Usage

```bash
neovector convert to-vector <input> [output] [flags]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `input`  | Path to the source image file or glob pattern |
| `output` | Output file path (omit for batch mode) |

## Flags

| Flag | Alias | Default | Description |
|------|-------|---------|-------------|
| `--format` | | `txt` | Output format: `txt`, `json`, `csv`, `bin` |
| `--grayscale` | `-g` | `false` | Convert to single-channel grayscale |
| `--channels` | | `rgb` | Select channels: `rgb`, `r`, `g`, `b`, `rg`, `rb`, `gb` |
| `--overwrite` | `-o` | `false` | Overwrite output file if it exists |
| `--verbose` | `-V` | `false` | Show detailed conversion progress |
| `--recursive` | `-r` | `false` | Process directories recursively (batch mode) |

## Examples

### Basic conversion (txt format)

```bash
neovector convert to-vector photo.jpg photo.txt
```

### JSON format

```bash
neovector convert to-vector photo.jpg photo.json
```

### CSV format (one pixel per row)

```bash
neovector convert to-vector photo.jpg photo.csv
```

### Binary format (compact, fast)

```bash
neovector convert to-vector photo.jpg photo.bin
```

### Grayscale

```bash
neovector convert to-vector photo.jpg gray.txt --grayscale
```

### Single channel

```bash
neovector convert to-vector photo.jpg red.txt --channels r
neovector convert to-vector photo.jpg rg.txt --channels rg
```

### Batch convert all PNGs in a directory

```bash
neovector convert to-vector "images/*.png" ./vectors --recursive
```

### Batch with verbose progress

```bash
neovector convert to-vector "*.jpg" ./vectors --recursive --verbose
```

## File Format

### Header

All text-based formats (txt, csv) include a header line:

```
# width=1920 height=1080 channels=3
```

This allows `to-image` to reconstruct the image without providing
width and height manually.

### TXT

```
# width=1920 height=1080 channels=3
255
128
64
...
```

### CSV

```
# width=1920 height=1080 channels=3
255,128,64
32,16,8
...
```

### JSON

Pretty-printed JSON array (no header):

```json
[
  255,
  128,
  64
]
```

### BIN

Raw little-endian uint8 triplets — most compact format.
No header; dimensions must be provided to `to-image`.

## Notes

- Input images are converted to RGB; alpha channel is discarded
- The vector size equals `width * height * channels`
- `--grayscale` and `--channels` cannot be combined
- In batch mode, output is a directory; output files are named after inputs
