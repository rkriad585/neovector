# config — View or Edit neovector Configuration

Display, modify, or interactively edit the neovector configuration.

## Usage

```bash
neovector config
neovector config --format <txt|json>
neovector config --proxy <url>
neovector config --edit
```

## Flags

| Flag | Type | Description |
|------|------|-------------|
| `--format` | string | Set default output format (`txt` or `json`) |
| `--proxy`  | string | Set proxy URL for self-update |
| `--edit`   | bool   | Open interactive TUI form editor |

## Examples

### View current configuration

```bash
neovector config
```

Output:
```
━━ neovector Configuration ━━

  Config file:
    ~/.config/neostore/neovector/config.toml

  [general]
    default_format = "txt"

  [network]
    proxy = ""

  Run 'neovector config --edit' to modify interactively.
```

### Set format directly

```bash
neovector config --format json
```

### Set proxy directly

```bash
neovector config --proxy http://proxy.company.com:8080
```

### Clear proxy

```bash
neovector config --proxy ""
```

### Interactive editor

```bash
neovector config --edit
```

Opens a TUI form with:
- **Default Output Format** — select between TXT and JSON
- **Proxy URL** — text input for the proxy address

Changes are saved immediately when the form is submitted.
If no changes are made, a "No changes made." message is shown.

## Configuration File

**Location:** `~/.config/neostore/neovector/config.toml`

```toml
[general]
default_format = "txt"

[network]
proxy = ""
```

## Runtime Reloading

Changes take effect immediately for subsequent commands within the same
session. The in-memory config is updated atomically and persisted to disk.

## Notes

- The config file is created automatically on first run
- The `--edit` flag uses the [HUH](https://github.com/charmbracelet/huh) TUI form library
- Changes made via `--edit` that match the current values are detected and skipped
