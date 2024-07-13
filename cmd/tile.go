package cmd

import (
	"fmt"
	"github.com/WildEgor/img2leaftlet/internal/handlers"
	"github.com/spf13/cobra"
	"image"
	"os"
	"path"
	"path/filepath"
)

// handler default tile handler
var handler = handlers.NewTileHandler()

// tileCmd represents tile cmd
var tileCmd = &cobra.Command{
	Use:   "tile",
	Short: "Tile any image to tiles",
	Long:  `Tile any image to tiles. Use .jpg or .png.`,
	Run: func(cmd *cobra.Command, _ []string) {
		inputPath := "image.jpg"
		outputPath := "output"
		tileSize := 256

		cmd.Flags().StringVarP(&inputPath, "input", "in", "image.jpg", "Specify input image path")
		cmd.Flags().StringVarP(&outputPath, "output", "out", "output", "Specify output dir path")
		cmd.Flags().IntVarP(&tileSize, "size", "s", 256, "Specify tile's size")

		imgfile, err := os.Open(inputPath)
		defer imgfile.Close() //nolint:all
		if err != nil {
			fmt.Println("file not found!")
			os.Exit(1)
		}

		getwd, err := os.Getwd()
		if err != nil {
			return
		}

		defaultOutput := filepath.ToSlash(path.Join(getwd, outputPath))
		fmt.Println(defaultOutput)

		if stat, err := os.Stat(defaultOutput); err != nil && stat == nil {
			//nolint:gosec
			if err := os.Mkdir(defaultOutput, 0777); err != nil {
				fmt.Println("cannot create output dir!")
				os.Exit(1)
			}
		}

		img, _, err := image.Decode(imgfile)
		if err != nil {
			fmt.Println("failed to decode image!")
			os.Exit(1)
		}

		if err := handler.Handle(&handlers.Args{
			Image:  img,
			Output: defaultOutput,
			Size:   tileSize,
		}); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

//nolint:all
func init() {
	rootCmd.AddCommand(tileCmd)
}
