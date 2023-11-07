package lib

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/codethread/dotty/lib/fp"
)

type ImplicitConfig struct {
	Home           string
	ConfigLocation string
}

// GetImplicitConfig gathers config information based on user's system, these are assumed as defaults if not overridden
func GetImplicitConfig() ImplicitConfig {
	HOME, err := os.UserHomeDir()

	if err != nil {
		log.Fatal("Could not find a suitable HOME value")
	}

	ConfigLocation := filepath.Join(HOME, "PersonalConfigs", ".dottyignore")

	return ImplicitConfig{
		Home:           HOME,
		ConfigLocation: ConfigLocation,
	}
}

type SetupConfig struct {
	DryRun      bool
	Force       bool
	From        string
	To          string
	ignored     []Matcher
	gitignores  []string
	HistoryFile string
	ImplicitConfig
}

func BuildSetupConfig(flags Flags, implicitConfig ImplicitConfig, target DottyTarget) SetupConfig {
	e := ExpandHome(implicitConfig.Home)
	var matcher []Matcher

	ignores := parseDottyIgnore(implicitConfig.ConfigLocation)
	if flags.Ignores != nil {
		ignores = append(ignores, *flags.Ignores...)
	}

	for _, re := range ignores {
		ignore, err := regexp.Compile(re)
		if err != nil {
			panic(err)
		}

		matcher = append(matcher, ignore)
	}

	safeName := strings.ReplaceAll(target.From, "/", "__")

	historyFile := e(fmt.Sprintf("~/.config/dotty/%s", safeName))
	dottyDir := e("~/.config/dotty")

	if _, err := os.Stat(dottyDir); os.IsNotExist(err) {
		os.MkdirAll(dottyDir, 0700)
	}

	return SetupConfig{
		From: e(target.From),
		To:   e(target.To),
		gitignores: []string{
			e(fmt.Sprintf("%s/.gitignore_global", target.From)),
			e(fmt.Sprintf("%s/.gitignore", target.From)),
		},
		ignored:        matcher,
		HistoryFile:    historyFile,
		ImplicitConfig: implicitConfig,
	}
}

func parseDottyIgnore(path string) (ignores []string) {
	file, err := os.ReadFile(path)

	if err != nil {
		return
	}

	// TODO: should likely use some valiation here
	lines := strings.Split(string(file), "\n")

	return fp.Filter(lines, func(s string) bool {
		return s != ""
	})
}

func ExpandHome(home string) func(string) string {
	return func(path string) string {
		if path == "~" {
			return home
		} else if strings.HasPrefix(path, "~/") {
			return filepath.Join(home, path[2:])
		} else {

			return path
		}
	}
}

type DottyTargets struct {
	Targets []DottyTarget
}

type DottyTarget struct {
	From string
	To   string
}

func GetDottyEnv() DottyTargets {
	dotty := os.Getenv("DOTTY")
	if dotty == "" {
		fmt.Println("env var DOTTY not set")
		os.Exit(1)
	}

	dirs := strings.Split(dotty, ":")
	length := len(dirs)

	if length%2 != 0 {
		fmt.Printf("DOTTY env var contains uneven target pairs, value of length %d\n\n%s\n", length, dotty)
		os.Exit(1)
	}

	var targets DottyTargets

	for x := 0; x < length; {
		t := DottyTarget{From: dirs[x], To: dirs[x+1]}
		targets.Targets = append(targets.Targets, t)
		x += 2
	}

	return targets
}
