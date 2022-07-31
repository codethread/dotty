/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/codethread/dotty/lib"

	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Symlink files into HOME",
	Long:  `Symlink files from TARGET into HOME, according to config`,
	Run: func(cmd *cobra.Command, args []string) {

		config := lib.BuildSetupConfig(
			lib.Flags{
				DryRun: lib.DryRun,
			},
			lib.GetImplicitConfig(),
		)

		lib.Setup(config)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
