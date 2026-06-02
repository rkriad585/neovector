# self-uninstall — Remove neovector from Your System

Fully uninstall neovector, including the binary, config files, and history.

## Usage

```bash
neovector --selfuninstall
neovector --uninstall
neovector -u
```

## Flags

| Flag | Description |
|------|-------------|
| `--selfuninstall` | Full command form |
| `--uninstall` | Alias for `--selfuninstall` |
| `-u` | Short form |

## What It Does

1. Resolves the running binary location
2. Determines platform-specific removal strategy
3. Removes config directory and files
4. Deletes the binary (or schedules deletion on Windows)

## Platform Behavior

| Platform | Binary Inside Config Dir | Binary Outside Config Dir |
|----------|-------------------------|---------------------------|
| **Windows** | Removes all config files except running exe, launches deferred `.bat` to delete remaining files (waits 1s, then `rmdir /s /q`) | Removes config dir directly, launches deferred `.bat` to delete the exe |
| **Linux/macOS** | Removes entire config dir, deletes binary directly | Same |

After uninstall, instructions are printed for removing the `neostore/neovector/bin` entry from your PATH.

## Examples

```bash
# Full command
neovector --selfuninstall

# Alias
neovector --uninstall

# Short flag
neovector -u
```

## Notes

- On Windows, the binary cannot delete itself while running, so a temporary batch script handles delayed deletion
- On Unix, the binary is removed immediately (inode remains valid for the running process)
- Config directory is always fully removed
- If config was already deleted, the tool handles it gracefully
- Reinstall with: `irm https://raw.githubusercontent.com/rkriad585/neovector/main/installer.ps1 \| iex` (Windows) or `curl -fsSL https://raw.githubusercontent.com/rkriad585/neovector/main/installer.sh \| sh` (Unix)
