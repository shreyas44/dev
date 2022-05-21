/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/rodaine/table"
	"github.com/shreyas44/dev/db"
	"github.com/spf13/cobra"
)

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:          "ps",
	Short:        "List all running processes",
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		dev, err := getDev()
		if err != nil {
			return err
		}

		t := table.New("NAME", "PID", "SATUS")
		for _, process := range dev.DB().ProcessesList() {
			status := string(process.Status)
			if process.Status == db.ProcessStatusExited {
				status += fmt.Sprintf(" (%d)", process.ExitCode)
			}

			t.AddRow(process.Name, process.PID, status)
		}
		t.Print()

		return nil
	},
}

func stringPtr(s string) *string {
	return &s
}

func uintPtr(i uint) *uint {
	return &i
}

func init() {
	rootCmd.AddCommand(psCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// psCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// psCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
