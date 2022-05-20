package cmd

import (
	"github.com/shreyas44/dev/dev"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs [service-name]",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, services []string) {
		logger := dev.NewLogger(services...)
		logger.Watch()
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	// logsCmd.Flags().BoolP("follow", "f", false, "Follow logs")
}
