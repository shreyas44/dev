/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"github.com/shreyas44/dev/dev"
	"github.com/spf13/cobra"
)

var tableStyle = ansi.StyleConfig{
	Table: ansi.StyleTable{
		RowSeparator:    stringPtr(""),
		ColumnSeparator: stringPtr("  "),
		CenterSeparator: stringPtr(""),
	},
}

// psCmd represents the ps command
var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "List all running processes",
	Run: func(cmd *cobra.Command, args []string) {
		wd, _ := os.Getwd()
		devNixPath, _ := dev.GetDevNixPath(wd)
		processes := devNixPath.Processes()
		if len(processes) == 0 {
			fmt.Println("No services running")
		} else {
			r, _ := glamour.NewTermRenderer(glamour.WithStyles(tableStyle))
			out, _ := r.Render(processes.Markdown())
			fmt.Print(out)
		}
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
