package cmd

import (
	"os"

	"github.com/shreyas44/dev/dev"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Shutdown dev environment and stop processes",
	Run: func(cmd *cobra.Command, args []string) {
		wd, _ := os.Getwd()
		dev, _ := dev.Get(wd)
		dev.Stop()
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
