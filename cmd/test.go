/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/codethread/dotty/lib"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		doStuff()
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}

func doStuff() {
	println("dr", lib.DryRun)

}
