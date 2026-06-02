# Themes — Customize the Look of neovector

neovector includes 13 built-in color themes that change the terminal output
colors. You can switch themes at any time — changes take effect immediately.

## Quick Start

```bash
# See all themes with color previews
neovector config theme

# Switch to a theme by name
neovector config theme vivid_nightfall

# Browse themes interactively
neovector config theme --edit
```

## Available Themes

| # | Name | Description |
|---|------|-------------|
| 1 | `dark` | Dark Theme — muted blues on dark background |
| 2 | `light` | Light Theme — whites and grays on light background |
| 3 | `sunny_beach_day` | Sunny Beach Day (Default) — teal, gold, coral |
| 4 | `olive_garden_feast` | Olive Garden Feast — greens, cream, amber |
| 5 | `summer_ocean_breeze` | Summer Ocean Breeze — red, cyan, navy |
| 6 | `refreshing_summer_fun` | Refreshing Summer Fun — sky blue, orange, gold |
| 7 | `black_gold_elegance` | Black & Gold Elegance — black, white, gold |
| 8 | `vibrant_color_fiesta` | Vibrant Color Fiesta — neon rainbow palette |
| 9 | `light_steel` | Light Steel — 9-shade gray monochrome |
| 10 | `golden_twilight` | Golden Twilight — deep navy with gold |
| 11 | `deep_sea` | Deep Sea — ocean blues and slate |
| 12 | `bright_green` | Bright Green — 8 shades of green |
| 13 | `vivid_nightfall` | Vivid Nightfall — deep purple to lavender |

## How Themes Work

Each theme defines 5+ hex colors that are mapped to semantic roles:

| Role | Used For |
|------|----------|
| **Primary** | Headings, labels, file paths |
| **Success** | Success messages, confirmations |
| **Warning** | Warnings, informational notes |
| **Error** | Error messages |
| **Accent** | Special highlights, decorative elements |

Themes with more than 5 colors use the extras for richer palette previews.

## Dark / Light Mode

The `[theme] mode` setting (in config.toml) selects the overall contrast
level. This is planned for future use to automatically adjust brightness
levels for better readability on dark vs light terminal backgrounds.

For now, choose a theme that matches your terminal:
- **Dark terminal** → `dark`, `golden_twilight`, `deep_sea`, `vivid_nightfall`
- **Light terminal** → `light`, `light_steel`, `black_gold_elegance`
- **Any terminal** → `sunny_beach_day`, `summer_ocean_breeze`, `refreshing_summer_fun`

## Interactive Theme Picker

The `--edit` flag opens a HUH TUI form that lets you browse themes with
live color swatches:

```bash
neovector config theme --edit
```

Each theme is shown with its name, a label, and a color swatch (████ bars
showing the actual palette). Select one and press Enter to apply it.

## Configuration

Theme settings are stored in `config.toml`:

```toml
[theme]
name = "sunny_beach_day"
mode = "dark"
```

- **name** — one of the 13 theme names listed above
- **mode** — `"dark"` or `"light"` (reserved for future brightness adaptation)

## Adding Custom Themes

To add a new theme, edit `internal/theme/theme.go` and add a new entry
to the `Themes` slice:

```go
{
    Name: "my_custom_theme", Label: "My Custom Theme",
    Colors: []string{"#hex1", "#hex2", "#hex3", "#hex4", "#hex5"},
},
```

The first color maps to **Primary**, the second to **Success**, and so on.
At least 5 colors are recommended for full coverage.
