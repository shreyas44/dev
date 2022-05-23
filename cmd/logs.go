package cmd

import (
	"github.com/shreyas44/dev/dev"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:          "logs [SERVICES]",
	Short:        "Show logs of services. Outputs logs of all running services if no services are specified",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, services []string) error {
		d, err := getDev()
		if err != nil {
			return err
		}

		logger, err := dev.NewLogger(d, services...)
		if err != nil {
			return err
		}

		logger.Watch()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)
	// logsCmd.Flags().BoolP("follow", "f", false, "Follow logs")
}
