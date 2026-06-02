# to-vector — Convert Image to Vector

Convert a raster image into a flat list of RGB pixel values.

## Usage

```bash
neovector convert to-vector <input> <output> [--format txt|json]
```

## Arguments

| Argument | Description |
|----------|-------------|
| `input`  | Path to the source image file |
| `output` | Path for the output vector file |

## Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `txt` | Output format: `txt` or `json` |

## Examples

### Basic conversion (txt format)

```bash
neovector convert to-vector photo.jpg photo.txt
```

### JSON format

```bash
neovector convert to-vector photo.jpg photo.json --format json
```

### Absolute output path

```bash
neovector convert to-vector photo.jpg /tmp/vectors/photo.txt
```

### Using an input from the output directory

```bash
neovector convert to-vector my_image.png my_vector.txt
```

If the output is a bare filename (no path), it is saved to
`~/Downloads/neostore/neovector/`.

## Output Format

### TXT

One integer per line:

```
255
128
64
...
```

### JSON

Pretty-printed JSON array:

```json
[
  255,
  128,
  64
]
```

## Notes

- Input images are converted to RGB; alpha channel is discarded
- The vector size equals `width * height * 3`
