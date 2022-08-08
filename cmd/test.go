/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

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
	value := 4
	pointer := &value

	*pointer = *pointer + 1

	fmt.Println(pointer, *pointer, value)

	mutate(&value)

	fmt.Println(value)

}

func mutate(i *int) {
	*i = *i * 3
}

type Foo struct {
	name string
}
