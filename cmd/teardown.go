/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/codethread/dotty/lib"
	"github.com/spf13/cobra"
)

// teardownCmd represents the teardown command
var teardownCmd = &cobra.Command{
	Use:   "teardown",
	Short: "Remove all files originally created by dotty",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config := lib.BuildSetupConfig(
			lib.Flags{
				DryRun: lib.DryRun,
			},
			lib.GetImplicitConfig(),
		)

		lib.Teardown(config)
	},
}

func init() {
	rootCmd.AddCommand(teardownCmd)
	teardownCmd.Flags().BoolVarP(&lib.DryRun, "dry-run", "d", false, "show the files that would be affected, without running actually changing anything")
}
