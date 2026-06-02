# Configuration

neovector stores its configuration at:

| Platform | Path |
|----------|------|
| Windows  | `%USERPROFILE%\.config\neostore\neovector\` |
| Linux    | `~/.config/neostore/neovector/` |
| macOS    | `~/.config/neostore/neovector/` |

## config.toml

Created automatically on first run.

```toml
[general]
default_format = "txt"

[network]
proxy = ""

[theme]
name = "sunny_beach_day"
mode = "dark"
```

### Options

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `general.default_format` | string | `"txt"` | Default vector format: `txt` or `json` |
| `network.proxy` | string | `""` | Proxy URL for self-update (e.g. `http://proxy:8080`) |
| `theme.name` | string | `"sunny_beach_day"` | Color theme name (see 13 built-in themes) |
| `theme.mode` | string | `"dark"` | Theme mode: `"dark"` or `"light"` |

## history.log

All conversions are logged with RFC3339 timestamps.

```
[2026-06-02T15:15:59+06:00] to-vector: photo.jpg -> vector.txt [txt, 4717548 values]
[2026-06-02T15:18:26+06:00] to-image: vector.txt -> restored.png [1254x1254, 4717548 values]
```

## Output Directory

Converted files (vectors, images) default to:

| Platform | Path |
|----------|------|
| Windows  | `%USERPROFILE%\Downloads\neostore\neovector\` |
| Linux    | `~/Downloads/neostore/neovector/` |
| macOS    | `~/Downloads/neostore/neovector/` |

Use an absolute or explicit relative path (e.g. `./output.txt`)
to save elsewhere.
