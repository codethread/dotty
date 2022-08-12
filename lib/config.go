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

	ignore, err := regexp.Compile("^_.*")
	ignore2, err := regexp.Compile("^.git$")
	ignore3, err := regexp.Compile("^node_modules$")
	if err != nil {
		panic(err)
	}

	matcher = append(matcher, ignore, ignore2, ignore3)

	return SetupConfig{
		DryRun:      true,
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
