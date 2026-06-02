package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/rkriad585/neovector/internal/config"
	"github.com/rkriad585/neovector/internal/log"
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert between images and vectors",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var toVectorCmd = &cobra.Command{
	Use:   "to-vector <input> <output>",
	Short: "Convert an image to a numerical vector",
	Long: `Convert an image file into a flat list of RGB pixel values.
Supports txt and json output formats.

Output is saved to ~/Downloads/neostore/neovector/ by default.
Provide an absolute or relative path to save elsewhere.

Example:
  neovector convert to-vector image.png vector.txt --format txt
  neovector convert to-vector image.jpg vector.json --format json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := config.ResolveInput(args[0])
		output := config.ResolveOutput(args[1])
		format, _ := cmd.Flags().GetString("format")

		Cyan.Printf("Converting %s to a numerical vector...\n", input)

		vector, err := imageToVector(input)
		if err != nil {
			Red.Printf("Error: %v\n", err)
			return err
		}

		if err := ensureOutputDir(); err != nil {
			return err
		}
		if err := saveVector(vector, output, format); err != nil {
			Red.Printf("Error saving vector file: %v\n", err)
			return err
		}

		Green.Printf("Success! Vector saved to %s (%s).\n", output, strings.ToUpper(format))
		fmt.Printf("  - Vector Size: %d values\n", len(vector))

		log.Write(fmt.Sprintf("to-vector: %s -> %s [%s, %d values]", input, output, format, len(vector)))
		return nil
	},
}

var toImageCmd = &cobra.Command{
	Use:   "to-image <input> <output> <width> <height>",
	Short: "Convert a vector back to an image",
	Long: `Reconstruct an image from a vector file.
You must provide the original width and height of the image.

Output is saved to ~/Downloads/neostore/neovector/ by default.

Example:
  neovector convert to-image vector.txt restored.png 1920 1080
  neovector convert to-image vector.json restored.png 800 600 --format json`,
	Args: cobra.ExactArgs(4),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := config.ResolveInput(args[0])
		output := config.ResolveOutput(args[1])
		format, _ := cmd.Flags().GetString("format")

		width, err := strconv.Atoi(args[2])
		if err != nil || width <= 0 {
			Red.Println("Error: width must be a positive integer")
			return fmt.Errorf("invalid width")
		}
		height, err := strconv.Atoi(args[3])
		if err != nil || height <= 0 {
			Red.Println("Error: height must be a positive integer")
			return fmt.Errorf("invalid height")
		}

		Cyan.Printf("Reconstructing image from %s...\n", input)

		vector, err := readVector(input, format)
		if err != nil {
			Red.Printf("Error reading vector file: %v\n", err)
			return err
		}

		if err := ensureOutputDir(); err != nil {
			return err
		}
		if err := vectorToImage(vector, width, height, output); err != nil {
			Red.Printf("Error creating image: %v\n", err)
			return err
		}

		Green.Printf("Success! Image saved to %s.\n", output)

		log.Write(fmt.Sprintf("to-image: %s -> %s [%dx%d, %d values]", input, output, width, height, len(vector)))
		return nil
	},
}

func ensureOutputDir() error {
	if err := config.EnsureOutputDir(); err != nil {
		Red.Printf("Error creating output directory: %v\n", err)
		return err
	}
	return nil
}

func imageToVector(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("image file not found at %s", path)
	}
	defer file.Close()

	src, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	bounds := src.Bounds()
	var data []int
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			data = append(data, int(r>>8), int(g>>8), int(b>>8))
		}
	}
	return data, nil
}

func vectorToImage(data []int, width, height int, outputPath string) error {
	expected := width * height * 3
	if len(data) != expected {
		return fmt.Errorf("vector size mismatch: expected %d values, got %d", expected, len(data))
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	idx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(data[idx]),
				G: uint8(data[idx+1]),
				B: uint8(data[idx+2]),
				A: 255,
			})
			idx += 3
		}
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	ext := strings.ToLower(filepath.Ext(outputPath))
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(out, img, &jpeg.Options{Quality: 95})
	default:
		return png.Encode(out, img)
	}
}

func saveVector(data []int, path, format string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch format {
	case "json":
		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		return encoder.Encode(data)
	default:
		w := bufio.NewWriter(file)
		for i, v := range data {
			if i > 0 {
				if err := w.WriteByte('\n'); err != nil {
					return err
				}
			}
			if _, err := fmt.Fprint(w, v); err != nil {
				return err
			}
		}
		return w.Flush()
	}
}

func readVector(path, format string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("vector file not found at %s", path)
	}
	defer file.Close()

	switch format {
	case "json":
		var data []int
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&data); err != nil {
			return nil, fmt.Errorf("invalid JSON vector file: %w", err)
		}
		return data, nil
	default:
		scanner := bufio.NewScanner(file)
		var data []int
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			v, err := strconv.Atoi(line)
			if err != nil {
				return nil, fmt.Errorf("invalid integer on line %d: %w", len(data)+1, err)
			}
			data = append(data, v)
		}
		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading vector file: %w", err)
		}
		return data, nil
	}
}

func init() {
	convertCmd.AddCommand(toVectorCmd)
	convertCmd.AddCommand(toImageCmd)

	toVectorCmd.Flags().String("format", "txt", "Output format (txt or json)")
	toImageCmd.Flags().String("format", "txt", "Input vector format (txt or json)")

	rootCmd.AddCommand(convertCmd)
}
