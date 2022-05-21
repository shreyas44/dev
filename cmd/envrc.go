/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"
)

// envrcCmd represents the envrc command
var envrcCmd = &cobra.Command{
	Use:   "envrc",
	Short: "Generate envrc file",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("PATH_add %s\n", path.Join(".dev-cli", "nix", "profile", "bin"))
	},
}

func init() {
	rootCmd.AddCommand(envrcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// envrcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// envrcCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
