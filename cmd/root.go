package cmd

import (
	"os"
	"path"

	"github.com/shreyas44/dev/dev"
	"github.com/spf13/cobra"
)

var filePath string

var rootCmd = &cobra.Command{
	Use:   "dev",
	Short: "Reproducable dev environments",
}

func getDev() (dev.Dev, error) {
	wd, err := os.Getwd()
	if filePath != "" {
		wd = path.Dir(filePath)
		dev.DevFileName = path.Base(filePath)
	}

	if err != nil {
		return dev.Dev{}, err
	}

	d, err := dev.Get(wd)
	if err != nil {
		return dev.Dev{}, err
	}

	return d, nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&filePath, "file", "f", "", "File to parse")
}
