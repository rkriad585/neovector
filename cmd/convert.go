package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rkriad585/neovector/internal/config"
	"github.com/rkriad585/neovector/internal/log"
)

type VectorData struct {
	Data   []int
	Width  int
	Height int
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert between images and vectors",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var (
	validFormats = []string{"txt", "json", "csv", "bin"}
)

func isValidFormat(f string) bool {
	for _, v := range validFormats {
		if f == v {
			return true
		}
	}
	return false
}

func resolveFormat(cmd *cobra.Command, filePath string) string {
	if cmd.Flags().Changed("format") {
		format, _ := cmd.Flags().GetString("format")
		return format
	}
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		return "json"
	case ".csv":
		return "csv"
	case ".bin":
		return "bin"
	}
	return config.GetDefaultFormat()
}

func changeExt(path, format string) string {
	ext := filepath.Ext(path)
	return path[:len(path)-len(ext)] + "." + format
}

func makeOutputPath(input, outputBase, format string, batch bool) string {
	if batch {
		return filepath.Join(outputBase, changeExt(filepath.Base(input), format))
	}
	return config.ResolveOutput(outputBase)
}

var toVectorCmd = &cobra.Command{
	Use:   "to-vector <input> [output]",
	Short: "Convert an image to a numerical vector",
	Long: `Convert an image file into a flat list of RGB pixel values.
Supports txt, json, csv and bin output formats.

The vector file includes a header with width, height, and channel count
so that to-image can reconstruct the image without manual dimensions.

When input is a glob (e.g. "*.png") and --recursive is set, all
matching files are converted in batch mode.

Examples:
  neovector convert to-vector image.png vector.txt
  neovector convert to-vector image.jpg vector.json --format json
  neovector convert to-vector image.png vector.csv --grayscale
  neovector convert to-vector image.png vector.bin --channels r
  neovector convert to-vector "*.png" ./output --recursive`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		grayscale, _ := cmd.Flags().GetBool("grayscale")
		channels, _ := cmd.Flags().GetString("channels")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		verbose, _ := cmd.Flags().GetBool("verbose")
		recursive, _ := cmd.Flags().GetBool("recursive")

		if grayscale && channels != "rgb" {
			return fmt.Errorf("--grayscale and --channels cannot be used together")
		}
		if !isValidChannelMode(channels) {
			return fmt.Errorf("invalid --channels %q: must be one of rgb, r, g, b, rg, rb, gb", channels)
		}

		inputs := expandInputs(args[0], recursive)
		if len(inputs) == 0 {
			return fmt.Errorf("no input files matched: %s", args[0])
		}

		batchMode := len(inputs) > 1

		format := getFormatFlag(cmd, "")
		if !batchMode && len(args) >= 2 && filepath.Ext(args[1]) != "" {
			format = getFormatFlag(cmd, args[1])
		}
		if !isValidFormat(format) {
			return fmt.Errorf("unsupported format %q: must be one of %s", format, strings.Join(validFormats, ", "))
		}

		outputBase := ""
		if batchMode {
			if len(args) == 2 {
				outputBase = args[1]
			} else {
				outputBase = config.GetOutputDir()
			}
		}

		ok, total := 0, len(inputs)
		var lastErr error
		for i, input := range inputs {
			output := ""
			if batchMode {
				output = makeOutputPath(input, outputBase, format, true)
			} else {
				output = makeOutputPath(input, args[1], format, false)
			}

			if !overwrite {
				if _, err := os.Stat(output); err == nil {
					Red.Printf("  Error [%d/%d]: %s already exists (use --overwrite)\n", i+1, total, output)
					lastErr = fmt.Errorf("output file exists")
					continue
				}
			}

			if verbose {
				Cyan.Printf("  [%d/%d] Reading: %s\n", i+1, total, input)
			}

			vd, err := imageToVector(input, grayscale, channels)
			if err != nil {
				Red.Printf("  Error [%d/%d] %s: %v\n", i+1, total, input, err)
				lastErr = err
				continue
			}

			if err := config.EnsureOutputDir(); err != nil {
				Red.Printf("  Error [%d/%d]: %v\n", i+1, total, err)
				lastErr = err
				continue
			}
			if verbose {
				Cyan.Printf("  [%d/%d] Writing: %s (%s, %d values)\n", i+1, total, output, strings.ToUpper(format), len(vd.Data))
			}
			if err := saveVector(vd, output, format); err != nil {
				Red.Printf("  Error [%d/%d] saving %s: %v\n", i+1, total, output, err)
				lastErr = err
				continue
			}

			Green.Printf("  [%d/%d] %s -> %s (%s, %dx%d)\n", i+1, total, input, output, strings.ToUpper(format), vd.Width, vd.Height)
			log.Write(fmt.Sprintf("to-vector: %s -> %s [%s, %d values]", input, output, format, len(vd.Data)))
			ok++
		}

		if batchMode {
			if lastErr != nil {
				Yellow.Printf("  Completed: %d/%d files converted (some failed)\n", ok, total)
				return lastErr
			}
			Green.Printf("  Completed: %d/%d files converted successfully.\n", ok, total)
		}
		return nil
	},
}

var toImageCmd = &cobra.Command{
	Use:   "to-image <input> <output> [width] [height]",
	Short: "Convert a vector back to an image",
	Long: `Reconstruct an image from a vector file.
Width and height can be omitted if the vector file has a metadata header
(supported by txt, csv formats). Without a header, both dimensions are required.

Supports txt, json, csv, and bin input formats.

When input is a glob (e.g. "*.txt") and --recursive is set, all
matching files are converted in batch mode (vectors must have headers
or use --width / --height flags for uniform dimensions).

Examples:
  neovector convert to-image vector.txt restored.png 1920 1080
  neovector convert to-image vector.json restored.jpg 800 600 --quality 90
  neovector convert to-image vector.bin restored.png
  neovector convert to-image "vectors/*.txt" ./images --recursive`,
	Args: cobra.RangeArgs(2, 4),
	RunE: func(cmd *cobra.Command, args []string) error {
		format := getFormatFlag(cmd, args[0])
		quality, _ := cmd.Flags().GetInt("quality")
		overwrite, _ := cmd.Flags().GetBool("overwrite")
		verbose, _ := cmd.Flags().GetBool("verbose")
		recursive, _ := cmd.Flags().GetBool("recursive")

		if !isValidFormat(format) {
			return fmt.Errorf("unsupported format %q: must be one of %s", format, strings.Join(validFormats, ", "))
		}
		if quality < 1 || quality > 100 {
			return fmt.Errorf("quality must be between 1 and 100, got %d", quality)
		}

		inputs := expandInputs(args[0], recursive)
		if len(inputs) == 0 {
			return fmt.Errorf("no input files matched: %s", args[0])
		}

		batchMode := len(inputs) > 1
		outputBase := args[1]
		sharedWidth, sharedHeight := 0, 0
		if len(args) >= 4 {
			w, err := strconv.Atoi(args[2])
			if err != nil || w <= 0 {
				return fmt.Errorf("invalid width: must be a positive integer")
			}
			h, err := strconv.Atoi(args[3])
			if err != nil || h <= 0 {
				return fmt.Errorf("invalid height: must be a positive integer")
			}
			if w > 65535 || h > 65535 {
				return fmt.Errorf("dimensions too large (maximum 65535 each)")
			}
			sharedWidth, sharedHeight = w, h
		}

		ok, total := 0, len(inputs)
		var lastErr error
		for i, input := range inputs {
			output := ""
			if batchMode {
				ext := ".png"
				if strings.HasSuffix(strings.ToLower(outputBase), ".jpg") || strings.HasSuffix(strings.ToLower(outputBase), ".jpeg") {
					ext = ".jpg"
				}
				output = filepath.Join(outputBase, changeExt(filepath.Base(input), ext))
				// strip the format extension first
				output = strings.TrimSuffix(output, "."+format)
				if !strings.HasSuffix(strings.ToLower(output), ext) {
					output += ext
				}
			} else {
				output = config.ResolveOutput(outputBase)
			}

			if !overwrite {
				if _, err := os.Stat(output); err == nil {
					Red.Printf("  Error [%d/%d]: %s already exists (use --overwrite)\n", i+1, total, output)
					lastErr = fmt.Errorf("output file exists")
					continue
				}
			}

			if verbose {
				Cyan.Printf("  [%d/%d] Reading: %s (%s)\n", i+1, total, input, strings.ToUpper(format))
			}

			vd, err := readVector(input, format)
			if err != nil {
				Red.Printf("  Error [%d/%d] %s: %v\n", i+1, total, input, err)
				lastErr = err
				continue
			}

			w, h := sharedWidth, sharedHeight
			if w == 0 && vd.Width > 0 {
				w, h = vd.Width, vd.Height
			}
			if w <= 0 || h <= 0 {
				Red.Printf("  Error [%d/%d] %s: dimensions not available; provide width and height as arguments\n", i+1, total, input)
				lastErr = fmt.Errorf("missing dimensions")
				continue
			}

			if err := config.EnsureOutputDir(); err != nil {
				Red.Printf("  Error [%d/%d]: %v\n", i+1, total, err)
				lastErr = err
				continue
			}
			if verbose {
				Cyan.Printf("  [%d/%d] Reconstructing %dx%d image...\n", i+1, total, w, h)
			}
			if err := vectorToImage(vd.Data, w, h, output, quality); err != nil {
				Red.Printf("  Error [%d/%d] creating %s: %v\n", i+1, total, output, err)
				lastErr = err
				continue
			}

			Green.Printf("  [%d/%d] %s -> %s (%dx%d)\n", i+1, total, input, output, w, h)
			log.Write(fmt.Sprintf("to-image: %s -> %s [%dx%d, %d values]", input, output, w, h, len(vd.Data)))
			ok++
		}

		if batchMode {
			if lastErr != nil {
				Yellow.Printf("  Completed: %d/%d files converted (some failed)\n", ok, total)
				return lastErr
			}
			Green.Printf("  Completed: %d/%d files converted successfully.\n", ok, total)
		}
		return nil
	},
}

func getFormatFlag(cmd *cobra.Command, filePath string) string {
	if cmd.Flags().Changed("format") {
		f, _ := cmd.Flags().GetString("format")
		return f
	}
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		return "json"
	case ".csv":
		return "csv"
	case ".bin":
		return "bin"
	}
	return config.GetDefaultFormat()
}

func expandInputs(pattern string, recursive bool) []string {
	// exact file?
	if fi, err := os.Stat(pattern); err == nil && !fi.IsDir() {
		return []string{pattern}
	}

	// directory with recursive?
	if fi, err := os.Stat(pattern); err == nil && fi.IsDir() && recursive {
		var files []string
		filepath.Walk(pattern, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			switch ext {
			case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".webp":
				files = append(files, path)
			}
			return nil
		})
		return files
	}

	// glob pattern
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return nil
	}
	return matches
}

func isValidChannelMode(channels string) bool {
	switch channels {
	case "rgb", "r", "g", "b", "rg", "rb", "gb":
		return true
	}
	return false
}

func clamp(v int) uint8 {
	if v < 0 {
		return 0
	}
	if v > 255 {
		return 255
	}
	return uint8(v)
}

func imageToVector(path string, grayscale bool, channels string) (VectorData, error) {
	file, err := os.Open(path)
	if err != nil {
		return VectorData{}, fmt.Errorf("failed to open image %s: %w", path, err)
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return VectorData{}, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := src.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	var data []int

	if grayscale {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := src.At(x, y).RGBA()
				gray := int((299*float64(r>>8) + 587*float64(g>>8) + 114*float64(b>>8)) / 1000)
				data = append(data, gray)
			}
		}
		return VectorData{Data: data, Width: w, Height: h}, nil
	}

	channelSel := parseChannels(channels)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			vr, vg, vb := int(r>>8), int(g>>8), int(b>>8)
			for _, ch := range channelSel {
				switch ch {
				case 'r':
					data = append(data, vr)
				case 'g':
					data = append(data, vg)
				case 'b':
					data = append(data, vb)
				}
			}
		}
	}
	return VectorData{Data: data, Width: w, Height: h}, nil
}

func parseChannels(channels string) []byte {
	if channels == "rgb" {
		return []byte{'r', 'g', 'b'}
	}
	return []byte(channels)
}

func vectorToImage(data []int, width, height int, outputPath string, quality int) error {
	expected := width * height * 3
	if len(data) != expected {
		return fmt.Errorf("vector size mismatch: expected %d values (3 per pixel), got %d", expected, len(data))
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	idx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: clamp(data[idx]),
				G: clamp(data[idx+1]),
				B: clamp(data[idx+2]),
				A: 255,
			})
			idx += 3
		}
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(outputPath))
	encodeErr := error(nil)
	switch ext {
	case ".jpg", ".jpeg":
		encodeErr = jpeg.Encode(out, img, &jpeg.Options{Quality: quality})
	default:
		encodeErr = png.Encode(out, img)
	}
	if closeErr := out.Close(); closeErr != nil && encodeErr == nil {
		encodeErr = fmt.Errorf("failed to close output file: %w", closeErr)
	}
	if encodeErr != nil {
		os.Remove(outputPath)
		return encodeErr
	}
	return nil
}

func headerLine(w, h int) string {
	return fmt.Sprintf("# width=%d height=%d channels=3", w, h)
}

func parseHeaderLine(line string) (width, height int, ok bool) {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "#") {
		return 0, 0, false
	}
	var w, h, c int
	n, _ := fmt.Sscanf(line, "# width=%d height=%d channels=%d", &w, &h, &c)
	if n == 3 && w > 0 && h > 0 && c > 0 {
		return w, h, true
	}
	return 0, 0, false
}

func saveVector(vd VectorData, path, format string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create vector file: %w", err)
	}

	encodeErr := error(nil)
	switch format {
	case "json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		encodeErr = encoder.Encode(vd.Data)
	case "csv":
		w := bufio.NewWriter(file)
		if _, err := fmt.Fprintln(w, headerLine(vd.Width, vd.Height)); err != nil {
			encodeErr = err
			break
		}
		for i, v := range vd.Data {
			if i > 0 {
				if (i)%3 == 0 {
					if _, err := w.WriteString("\n"); err != nil {
						encodeErr = err
						break
					}
				} else {
					if err := w.WriteByte(','); err != nil {
						encodeErr = err
						break
					}
				}
			}
			if encodeErr != nil {
				break
			}
			if _, err := fmt.Fprint(w, v); err != nil {
				encodeErr = err
				break
			}
		}
		if encodeErr == nil && len(vd.Data) > 0 {
			if _, err := w.WriteString("\n"); err != nil {
				encodeErr = err
			}
		}
		if encodeErr == nil {
			encodeErr = w.Flush()
		}
	case "bin":
		var buf bytes.Buffer
		buf.Grow(len(vd.Data))
		for _, v := range vd.Data {
			if v < 0 || v > 255 {
				return fmt.Errorf("value %d out of range for binary format (0-255)", v)
			}
			buf.WriteByte(byte(v))
		}
		_, encodeErr = io.Copy(file, &buf)
	default:
		w := bufio.NewWriter(file)
		if _, err := fmt.Fprintln(w, headerLine(vd.Width, vd.Height)); err != nil {
			encodeErr = err
			break
		}
		for i, v := range vd.Data {
			if i > 0 {
				if err := w.WriteByte('\n'); err != nil {
					encodeErr = err
					break
				}
			}
			if _, err := fmt.Fprint(w, v); err != nil {
				encodeErr = err
				break
			}
		}
		if encodeErr == nil {
			encodeErr = w.Flush()
		}
	}
	if closeErr := file.Close(); closeErr != nil && encodeErr == nil {
		encodeErr = fmt.Errorf("failed to close vector file: %w", closeErr)
	}
	if encodeErr != nil {
		os.Remove(path)
		return encodeErr
	}
	return nil
}

func readVector(path, format string) (VectorData, error) {
	file, err := os.Open(path)
	if err != nil {
		return VectorData{}, fmt.Errorf("failed to open vector file %s: %w", path, err)
	}
	defer file.Close()

	switch format {
	case "json":
		var data []int
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&data); err != nil {
			return VectorData{}, fmt.Errorf("invalid JSON vector file: %w", err)
		}
		return VectorData{Data: data}, nil
	case "csv":
		raw, err := io.ReadAll(file)
		if err != nil {
			return VectorData{}, fmt.Errorf("failed to read CSV vector file: %w", err)
		}
		var w, h int
		body := raw
		if idx := bytes.IndexByte(raw, '\n'); idx >= 0 {
			firstLine := string(raw[:idx])
			if w2, h2, ok := parseHeaderLine(firstLine); ok {
				w, h = w2, h2
				body = raw[idx+1:]
			}
		}
		body = bytes.ReplaceAll(body, []byte{'\r'}, nil)
		var data []int
		for _, field := range bytes.Split(body, []byte{','}) {
			field = bytes.TrimSpace(field)
			if len(field) == 0 {
				continue
			}
			for _, part := range bytes.Split(field, []byte{'\n'}) {
				part = bytes.TrimSpace(part)
				if len(part) == 0 {
					continue
				}
				v, err := strconv.Atoi(string(part))
				if err != nil {
					return VectorData{}, fmt.Errorf("invalid integer in CSV: %q: %w", string(part), err)
				}
				data = append(data, v)
			}
		}
		return VectorData{Data: data, Width: w, Height: h}, nil
	case "bin":
		raw, err := io.ReadAll(file)
		if err != nil {
			return VectorData{}, fmt.Errorf("failed to read binary vector file: %w", err)
		}
		data := make([]int, len(raw))
		for i, b := range raw {
			data[i] = int(b)
		}
		return VectorData{Data: data}, nil
	default:
		reader := bufio.NewReader(file)
		var w, h int
		// peek at first line for header
		firstLine, err := reader.ReadString('\n')
		if err == nil {
			if w2, h2, ok := parseHeaderLine(firstLine); ok {
				w, h = w2, h2
			} else {
				// first line is data, re-process from beginning
				reader = bufio.NewReader(io.MultiReader(
					strings.NewReader(firstLine),
					file,
				))
			}
		} else {
			// empty or single-line file without newline
			firstLine = strings.TrimSpace(firstLine)
			if firstLine != "" {
				if w2, h2, ok := parseHeaderLine(firstLine); ok {
					w, h = w2, h2
				} else {
					// it's data, re-process
					reader = bufio.NewReader(io.MultiReader(
						strings.NewReader(firstLine+"\n"),
						file,
					))
				}
			}
		}

		scanner := bufio.NewScanner(reader)
		var data []int
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			v, err := strconv.Atoi(line)
			if err != nil {
				return VectorData{}, fmt.Errorf("invalid integer on line %d: %w", len(data)+1, err)
			}
			data = append(data, v)
		}
		if err := scanner.Err(); err != nil {
			return VectorData{}, fmt.Errorf("error reading vector file: %w", err)
		}
		return VectorData{Data: data, Width: w, Height: h}, nil
	}
}

func init() {
	convertCmd.AddCommand(toVectorCmd)
	convertCmd.AddCommand(toImageCmd)

	toVectorCmd.Flags().String("format", "txt", "Output format (txt, json, csv, bin)")
	toVectorCmd.Flags().BoolP("grayscale", "g", false, "Convert to single-channel grayscale")
	toVectorCmd.Flags().String("channels", "rgb", "Select channels: rgb, r, g, b, rg, rb, gb")
	toVectorCmd.Flags().BoolP("overwrite", "o", false, "Overwrite output file if it exists")
	toVectorCmd.Flags().BoolP("verbose", "V", false, "Show detailed conversion progress")
	toVectorCmd.Flags().BoolP("recursive", "r", false, "Process directories recursively (batch mode)")

	toImageCmd.Flags().String("format", "txt", "Input vector format (txt, json, csv, bin)")
	toImageCmd.Flags().Int("quality", 95, "JPEG output quality (1-100)")
	toImageCmd.Flags().BoolP("overwrite", "o", false, "Overwrite output file if it exists")
	toImageCmd.Flags().BoolP("verbose", "V", false, "Show detailed conversion progress")
	toImageCmd.Flags().BoolP("recursive", "r", false, "Process directories recursively (batch mode)")

	rootCmd.AddCommand(convertCmd)
}
