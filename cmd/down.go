package cmd

import (
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:          "down",
	Short:        "Shutdown dev environment and stop processes",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dev, err := getDev()
		if err != nil {
			return err
		}

		dev.Stop()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
