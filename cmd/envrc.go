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
		fmt.Printf("PATH_add %s\n", path.Join(".dev", "nix", "profile", "bin"))
	},
}

func init() {
	rootCmd.AddCommand(envrcCmd)
}
