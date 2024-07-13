// Package cmd ...
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "i2l",
	Short: "",
	Long:  ``,
}

// Execute handle cmd
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init default image formats and call cobra init
//
//nolint:all
func init() {
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)

	cobra.OnInitialize()
}
