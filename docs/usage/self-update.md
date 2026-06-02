# self-update — Update neovector to the Latest Version

Check for and apply updates to the neovector binary from GitHub.

## Usage

```bash
neovector self-update [--proxy <url>]
```

## Flags

| Flag | Alias | Description |
|------|-------|-------------|
| `--proxy` | `-p` | Proxy URL for the update (e.g. `http://proxy:8080`) |

## Examples

### Check and update

```bash
neovector self-update
```

### Use a proxy

```bash
neovector self-update --proxy http://proxy.company.com:8080
```

## What It Does

1. Fetches the latest version from `raw.githubusercontent.com/rkriad585/neovector/main/.version`
2. Compares it with the current running version
3. If up-to-date: prints a success message and exits
4. If newer: detects OS/arch and downloads the matching release binary from GitHub
5. Replaces the current binary safely:

| Platform | Strategy |
|----------|----------|
| Windows  | Renames current exe → `.old`, places new binary, keeps rollback |
| Linux/macOS | Renames new binary over current exe, applies `chmod +x` |

## Proxy

The proxy can also be set persistently in `config.toml`:

```toml
[network]
proxy = "http://proxy.company.com:8080"
```

The `--proxy` flag takes precedence over the config file.

## Notes

- A download progress counter shows MB transferred
- On Windows the old binary is kept as `neovector.exe.old` for rollback
- Requires a working internet connection to GitHub
- The update downloads from `github.com/rkriad585/neovector/releases`
