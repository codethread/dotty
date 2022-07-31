package lib

import (
	"fmt"
	"github.com/codethread/dotty/lib/setup"
	"regexp"
)

func Setup(config SetupConfig) {
	// teardown
	// get ignored patterns
	setup.GetAllLinkableFiles(config.From, &config.Ignored)
	// sort by length to keep directories at the end for cleanup
	// for each file, link it into home
	// for all success, add them to a teardown file
	// for all failures, present a warning in the console

	fmt.Println("setup called")

}

type SetupConfig struct {
	DryRun      bool
	From        string
	To          string
	Ignored     []regexp.Regexp
	HistoryFile string
}

func BuildSetupConfig(flags Flags, implicitConfig ImplicitConfig) SetupConfig {
	panic("not implemented")
}
