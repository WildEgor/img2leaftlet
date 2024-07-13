// Package cmd ...
package cmd

import (
	"github.com/spf13/cobra"
	"image"
	"image/jpeg"
	"image/png"
	"log/slog"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "i2l",
	Short: "",
	Long:  ``,
}

// Execute default image formats and call cobra init
func Execute() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	cobra.OnInitialize()

	InitTileCmd()

	if err := rootCmd.Execute(); err != nil {
		slog.Error("error execute: ", slog.Any("err", err))
		os.Exit(1)
	}
}
