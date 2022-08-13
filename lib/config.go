package lib

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
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

	ConfigLocation := path.Join(HOME, ".config", "dotty")

	return ImplicitConfig{
		Home:           HOME,
		ConfigLocation: ConfigLocation,
	}
}

type SetupConfig struct {
	DryRun      bool
	From        string
	To          string
	ignored     []Matcher
	gitignores  []string
	HistoryFile string
}

func BuildSetupConfig(flags Flags, implicitConfig ImplicitConfig) SetupConfig {
	e := expand(implicitConfig.Home)
	var matcher []Matcher

	ignores := []string{
		"^_.*",
		"^.git$",
		"^node_modules$",
		".gitignore",
		"README",
	}

	for _, re := range ignores {
		ignore, err := regexp.Compile(re)
		if err != nil {
			panic(err)
		}

		matcher = append(matcher, ignore)
	}

	return SetupConfig{
		DryRun:      flags.DryRun,
		From:        e("~/PersonalConfigs"),
		To:          e("~"),
		gitignores:  []string{e("~/PersonalConfigs/.gitignore_global"), e("~/PersonalConfigs/.gitignore")},
		ignored:     matcher,
		HistoryFile: e("~/.dotty"),
	}
}

func expand(home string) func(string) string {
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
