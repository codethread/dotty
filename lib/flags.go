package lib

// Flag cli arg
var DryRun bool
var Ignores []string

type Flags struct {
	DryRun  bool
	Ignores *[]string
}
