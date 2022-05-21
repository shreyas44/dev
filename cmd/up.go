package cmd

import (
	"os"

	"github.com/shreyas44/dev/dev"

	"github.com/spf13/cobra"
)

// runs "dev activate" before starting services
var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Initialize dev environment and start processes",
	Run: func(cmd *cobra.Command, args []string) {
		wd, _ := os.Getwd()
		dev, _ := dev.Get(wd)
		dev.Init()
		dev.Stop()
		dev.Start()
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.PersistentFlags().BoolP("detached", "d", false, "Run in detached mode")
	upCmd.PersistentFlags().BoolP("nix-env", "e", false, "Run in detached mode")
}
