/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Dotty",
	Long:  `All software has versions. This is Dotty's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.0.10")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
