# PixelVector

A powerful, TUI-driven command-line tool for seamless conversion between images and their numerical vector representations.

![PixelVector Demo](https://storage.googleapis.com/gemini-assistant-tool-files/2024/05/pixel-vector-demo.gif)  <!-- Replace with an actual demo GIF if you create one -->

---

## ✨ Key Features

- **All-in-One Tool:** A single script to handle all conversions and checks.
- **Rich TUI:** A beautiful and intuitive interface powered by the `rich` library, featuring styled panels, tables, and progress bars.
- **Clear Sub-commands:** Easy-to-use sub-commands (`convert`, `check`) for clear and separate actions.
- **Multiple Formats:** Supports both plain text (`.txt`) and JSON (`.json`) for vector files.
- **Dimension Analysis:** Automatically calculates and displays possible image dimensions from a vector file.

---

## 🚀 Installation

1.  **Clone the repository (or download the script):**
    ```bash
    git clone https://github.com/your-username/pixelvector.git
    cd pixelvector
    ```

2.  **Install the required Python libraries:**
    PixelVector requires `Pillow` for image processing and `rich` for the enhanced TUI. You can install them using pip:
    ```bash
    pip install Pillow rich
    ```

---

## 🛠️ Usage

The tool is organized into two main sub-commands: `convert` and `check`.

### 1. Convert an Image to a Vector

This command takes an image file and converts it into a numerical vector.

```bash
python img-vec-tool.py convert to-vector <input_image> <output_file> [--format <format>]
```

**Example:**
```bash
python img-vec-tool.py convert to-vector my_image.png vector.txt --format txt
```

### 2. Convert a Vector to an Image

This command reconstructs an image from a vector file. You must provide the correct width and height of the target image.

```bash
python img-vec-tool.py convert to-image <input_file> <output_image> <width> <height> [--format <format>]
```

**Example:**
```bash
python img-vec-tool.py convert to-image vector.txt new_image.png 1920 1080
```

### 3. Check Dimensions

This command is used to inspect the dimensions of an image or the size and possible dimensions of a vector file. This is very useful for finding the correct `width` and `height` for the `to-image` command.

```bash
python img-vec-tool.py check [--image <path>] [--vector <path>]
```

**Examples:**

- **Check an image:**
  ```bash
  python img-vec-tool.py check --image my_image.png
  ```

- **Check a vector file:**
  ```bash
  python img-vec-tool.py check --vector vector.txt
  ```

- **Check both at once for a quick comparison:**
  ```bash
  python img-vec-tool.py check --image my_image.png --vector vector.txt
  ```

---

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
