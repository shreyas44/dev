package cmd

import (
	"github.com/spf13/cobra"
)

// runs "dev activate" before starting services
var upCmd = &cobra.Command{
	Use:          "up",
	Short:        "Initialize dev environment and start services",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dev, err := getDev()
		if err != nil {
			return err
		}

		dev.Init()
		dev.Stop()
		dev.Start()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
	upCmd.PersistentFlags().BoolP("detached", "d", false, "Run in detached mode")
	upCmd.PersistentFlags().BoolP("nix-env", "e", false, "Run in detached mode")
}
