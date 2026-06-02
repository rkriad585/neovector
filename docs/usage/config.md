# config — View or Edit neovector Configuration

Display, modify, or interactively edit the neovector configuration.

## Usage

```bash
neovector config
neovector config --format <txt|json>
neovector config --proxy <url>
neovector config --edit
neovector config theme                  # Show theme info
neovector config theme <name>           # Switch theme
neovector config theme --edit           # Interactive theme picker
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

  [theme]
    name = "sunny_beach_day"
    mode = "dark"

  Run 'neovector config --edit' or 'neovector config theme --edit' to modify.
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
- **Color Theme** — select from 13 available themes

Changes are saved immediately when the form is submitted.
If no changes are made, a "No changes made." message is shown.

## Theme Management

See the [Themes](themes.md) page for full details on all available themes
and how to switch between them.

## Configuration File

**Location:** `~/.config/neostore/neovector/config.toml`

```toml
[general]
default_format = "txt"

[network]
proxy = ""

[theme]
name = "sunny_beach_day"
mode = "dark"
```

## Runtime Reloading

Changes take effect immediately for subsequent commands within the same
session. The in-memory config is updated atomically and persisted to disk.
Theme changes are applied immediately to the terminal output.

## Notes

- The config file is created automatically on first run
- The `--edit` flag uses the [HUH](https://github.com/charmbracelet/huh) TUI form library
- Changes made via `--edit` that match the current values are detected and skipped
