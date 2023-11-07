package lib

// Flag cli arg
var DryRun bool
var Ignores []string
var Force bool

type Flags struct {
	DryRun  bool
	Ignores *[]string
	Force   bool
}
