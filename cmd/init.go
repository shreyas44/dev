/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
)

// initCmd represents the activate command
var initCmd = &cobra.Command{
	Use:          "init",
	Short:        "Installs dependencies and sets up environment",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dev, err := getDev()
		if err != nil {
			return err
		}

		dev.Init()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// activateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// activateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().BoolP("pure", "p", false, "Overwrite PATH completely except nix-env")
}
