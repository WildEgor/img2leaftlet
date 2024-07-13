package cmd

import (
	"github.com/WildEgor/img2leaftlet/internal/handlers"
	"github.com/spf13/cobra"
	"image"
	"log/slog"
	"os"
	"path"
	"path/filepath"
)

// handler default tile handler
var handler = handlers.NewTileHandler()

var (
	inputPath  = "image.jpg"
	outputPath = "tiles"
	tileSize   = 256
)

// tileCmd represents tile cmd
var tileCmd = &cobra.Command{
	Use:   "tile",
	Short: "Tile any image to tiles",
	Long:  `Tile any image to tiles. Use .jpg or .png.`,
	Run: func(_ *cobra.Command, _ []string) {
		imgfile, err := os.Open(inputPath)
		defer imgfile.Close() //nolint:all // ...
		if err != nil {
			slog.Error("fail open image file!", slog.Any("err", err))
			os.Exit(1)
		}

		getwd, err := os.Getwd()
		if err != nil {
			return
		}

		defaultOutput := filepath.ToSlash(path.Join(getwd, outputPath))
		slog.Debug("output path: ", slog.Any("value", defaultOutput))

		if stat, err := os.Stat(defaultOutput); err != nil && stat == nil {
			//nolint:gosec // ...
			if err := os.Mkdir(defaultOutput, 0777); err != nil {
				slog.Error("read output dir fail!", slog.Any("err", err))
				os.Exit(1)
			}
		}

		img, _, err := image.Decode(imgfile)
		if err != nil {
			slog.Error("failed to decode image!", slog.Any("err", err))
			os.Exit(1)
		}

		if err := handler.Handle(&handlers.Args{
			Image:  img,
			Output: defaultOutput,
			Size:   tileSize,
		}); err != nil {
			slog.Error("fail tile!", slog.Any("err", err))
			os.Exit(1)
		}
	},
}

// InitTileCmd ...
func InitTileCmd() {
	tileCmd.Flags().StringVarP(&inputPath, "input", "i", "image.jpg", "Specify input image path")
	tileCmd.Flags().StringVarP(&outputPath, "output", "o", "output", "Specify output dir path")
	tileCmd.Flags().IntVarP(&tileSize, "size", "s", 256, "Specify tile's size")
	rootCmd.AddCommand(tileCmd)
}
