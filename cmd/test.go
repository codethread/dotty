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
	Use:   "test file1 file2 ...files",
	Short: "Test a file(s) to see if dotty would link this file",
	Long: `Pass one or more files to Dotty to see if those files
will get symlinked accoring to dotty's heuristic

Exits with a non-zero status if any file is not a valid dotty file

$ dotty test some-file.txt some-other-file.txt

Files starting with '/' or '~' will be treated as absoulete or
expanded to HOME respectively

All other files will be appended to dotty config FROM location.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		config := lib.BuildSetupConfig(
			lib.Flags{
				Ignores: &lib.Ignores,
			},
			lib.GetImplicitConfig(),
		)

		lib.TestFiles(config, args)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().StringSliceVarP(&lib.Ignores, "ignore", "i", []string{}, "regex ignore patterns, e.g -i \"foo*\" -i \".*bar$\"")
}
