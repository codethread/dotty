package lib

import (
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
	From        string
	To          string
	ignored     []Matcher
	gitignores  []string
	HistoryFile string
	ImplicitConfig
}

func BuildSetupConfig(flags Flags, implicitConfig ImplicitConfig) SetupConfig {
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

	return SetupConfig{
		DryRun:         flags.DryRun,
		From:           e("~/PersonalConfigs"),
		To:             e("~"),
		gitignores:     []string{e("~/PersonalConfigs/.gitignore_global"), e("~/PersonalConfigs/.gitignore")},
		ignored:        matcher,
		HistoryFile:    e("~/.dotty"),
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
