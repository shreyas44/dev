package cmd

import (
	"os"

	"github.com/shreyas44/dev/dev"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dev",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func getDev() (dev.Dev, error) {
	wd, err := os.Getwd()
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
	rootCmd.PersistentFlags().StringP("file", "f", "", "File to parse")
}
