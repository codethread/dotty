/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/codethread/dotty/lib"
	"github.com/spf13/cobra"
)

// teardownCmd represents the teardown command
var teardownCmd = &cobra.Command{
	Use:   "teardown",
	Short: "Remove all files originally created by dotty",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("teardown called")

		HOME, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		fS := os.DirFS(HOME)
		f, err := fs.ReadFile(fS, ".dotty")

		if err != nil {
			panic(err)
		}

		files := lib.FromGOB64(string(f))
		files.Walk(lib.Visitor{
			File: func(dir string, file string) { fmt.Println("->", dir, file) },
			Dir:  func(dir string) { fmt.Println("dd", dir) },
		})
	},
}

func init() {
	rootCmd.AddCommand(teardownCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// teardownCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// teardownCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
