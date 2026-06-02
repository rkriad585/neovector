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

# Convert the vector back to an image
neovector convert to-image vector.txt restored.png 1920 1080
```

## Contents

- [Setup Guides](setup/windows.md)
- [Usage Guides](usage/to-vector.md)
- [Config](usage/config.md)
- [Self-Update](usage/self-update.md)
- [Self-Uninstall](usage/self-uninstall.md)
- [Configuration](configuration/config.md)
- [Development](development/build.md)
- [Troubleshooting](troubleshooting.md)
