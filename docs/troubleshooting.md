# Troubleshooting

## "command not found" after install

**Windows:** Restart your terminal, or run:

```powershell
refreshenv
```

**Linux/macOS:** Restart your shell, or run:

```bash
source ~/.bashrc
# or
source ~/.zshrc
# or
hash -r
```

## "vector file not found"

neovector looks for input files in the current directory and the
default output directory (`~/Downloads/neostore/neovector/`).

If your file is in neither location, provide the full path:

```bash
neovector convert to-image /path/to/vector.txt output.png 800 600
```

## "vector size mismatch"

The vector length must exactly equal `width * height * 3`.

Use `neovector check --vector <file>` to see all valid dimension pairs.

## "image file not found"

Ensure the image path is correct. Use absolute paths if needed:

```bash
neovector convert to-vector C:\Users\me\photo.jpg vector.txt
```

## Banner shows "unknown" for commit

The commit hash is injected via build ldflags. Plain `go build` will show
"unknown". Use `make build` or `build.ps1`/`build.sh` to include the commit.

## Build fails on Windows

Ensure you have Go installed and `git` in your PATH:

```powershell
go version
git --version
```

## Permission denied (Linux/macOS)

Make the binary executable:

```bash
chmod +x neovector
```

## Config directory not created

The config directory is created on first command run. Run any command
to initialize it:

```bash
neovector --version
```
