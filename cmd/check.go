package cmd

import (
	"fmt"
	"image"
	"math"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check dimensions of images or vectors",
	Long: `Inspect the dimensions of an image file or the size and possible
dimensions of a vector file.

Examples:
  neovector check --image image.png
  neovector check --vector vector.txt
  neovector check --image image.png --vector vector.txt`,
	RunE: func(cmd *cobra.Command, args []string) error {
		imagePath, _ := cmd.Flags().GetString("image")
		vectorPath, _ := cmd.Flags().GetString("vector")
		format, _ := cmd.Flags().GetString("format")

		if imagePath == "" && vectorPath == "" {
			return cmd.Help()
		}

		Yellow.Println("━━ Checking Dimensions ━━")

		if imagePath != "" {
			file, err := os.Open(imagePath)
			if err != nil {
				Red.Printf("Error: image file not found at %s\n", imagePath)
				return err
			}
			defer file.Close()

			img, _, err := image.Decode(file)
			if err != nil {
				Red.Printf("Error: failed to decode image: %v\n", err)
				return err
			}

			bounds := img.Bounds()
			w := bounds.Dx()
			h := bounds.Dy()
			pixels := w * h

			fmt.Println()
			Cyan.Printf("  Image: %s\n", imagePath)
			fmt.Printf("  ├─ Dimensions:     %d x %d\n", w, h)
			fmt.Printf("  ├─ Total Pixels:   %d\n", pixels)
			fmt.Printf("  └─ RGB Vector Size: %d\n", pixels*3)
		}

		if vectorPath != "" {
			data, err := readVector(vectorPath, format)
			if err != nil {
				Red.Printf("Error: %v\n", err)
				return err
			}

			size := len(data)
			fmt.Println()
			Cyan.Printf("  Vector: %s\n", vectorPath)
			fmt.Printf("  ├─ Format:         %s\n", strings.ToUpper(format))
			fmt.Printf("  ├─ Vector Size:    %d values\n", size)
			fmt.Printf("  └─ Possible RGB Dimensions:\n")

			if size%3 == 0 {
				pixels := size / 3
				fmt.Printf("       (Total Pixels: %d)\n", pixels)
				limit := int(math.Sqrt(float64(pixels)))
				maxWidth := 0
				pairs := make([][2]int, 0)
				for i := 1; i <= limit; i++ {
					if pixels%i == 0 {
						w, h := i, pixels/i
						pairs = append(pairs, [2]int{w, h})
						if w != h {
							pairs = append(pairs, [2]int{h, w})
						}
						if w > maxWidth {
							maxWidth = w
						}
						if h > maxWidth {
							maxWidth = h
						}
					}
				}

				headerW := len(fmt.Sprintf("%d", maxWidth))
				if headerW < 5 {
					headerW = 5
				}
				headerH := len("Height")
				if headerH < len(fmt.Sprintf("%d", maxWidth)) {
					headerH = len(fmt.Sprintf("%d", maxWidth))
				}
				if headerH < 6 {
					headerH = 6
				}

				fmt.Printf("       ┌─ %s ┬─ %s ┐\n",
					strings.Repeat("─", headerW+2),
					strings.Repeat("─", headerH+2))
				fmt.Printf("       │ %*s │ %*s │\n",
					-headerW-2, fmt.Sprintf(" %s ", "Width"),
					-headerH-2, fmt.Sprintf(" %s ", "Height"))
				fmt.Printf("       ├─ %s ┼─ %s ┤\n",
					strings.Repeat("─", headerW+2),
					strings.Repeat("─", headerH+2))
				for _, p := range pairs {
					fmt.Printf("       │ %*d │ %*d │\n",
						-headerW-2, p[0],
						-headerH-2, p[1])
				}
				fmt.Printf("       └─ %s ┴─ %s ┘\n",
					strings.Repeat("─", headerW+2),
					strings.Repeat("─", headerH+2))
			} else {
				Yellow.Println("       (Vector size is not divisible by 3; cannot determine RGB dimensions.)")
			}
		}

		return nil
	},
}

func init() {
	checkCmd.Flags().StringP("image", "i", "", "Path to the image file")
	checkCmd.Flags().StringP("vector", "v", "", "Path to the vector file")
	checkCmd.Flags().String("format", "txt", "Vector file format (txt or json)")

	rootCmd.AddCommand(checkCmd)
}
