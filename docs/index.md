# neovector Documentation

neovector converts between raster images and their numerical vector
representations — flat lists of RGB pixel values. It is a single static
Go binary with no runtime dependencies.

## Quick Start

```bash
# Convert an image to a vector
neovector convert to-vector image.png vector.txt

# Check its dimensions
neovector check --image image.png

# Convert the vector back to an image (header auto-detects dimensions)
neovector convert to-image vector.txt restored.png

# Or provide dimensions manually
neovector convert to-image vector.txt restored.png 1920 1080

# Batch convert all PNGs
neovector convert to-vector "images/*.png" ./vectors --recursive
```

## Contents

- [Setup Guides](setup/windows.md)
- [Usage Guides](usage/to-vector.md)
- [Config](usage/config.md)
- [Themes](usage/themes.md)
- [Self-Update](usage/self-update.md)
- [Self-Uninstall](usage/self-uninstall.md)
- [Configuration](configuration/config.md)
- [Development](development/build.md)
- [Troubleshooting](troubleshooting.md)
