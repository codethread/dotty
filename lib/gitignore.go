package lib

import (
	"os"

	ignore "github.com/sabhiram/go-gitignore"

	"github.com/codethread/dotty/lib/fp"
)

type Matcher interface {
	MatchString(f string) bool
}

type Matchers []Matcher

func (ignores Matchers) Matches(path string) bool {
	for _, ignore := range ignores {
		if ignore.MatchString(path) {
			return true
		}
	}
	return false
}

func GetAllIgnoredPatterns(config SetupConfig) []Matcher {
	ignored := getIgnoredPatterns(config.gitignores)
	return append(ignored, config.ignored...)
}

func getIgnoredPatterns(ignoreFiles []string) []Matcher {
	existingFiles := fp.Filter(ignoreFiles, func(file string) bool {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return false
		}
		return true
	})

	return fp.PromiseAll(existingFiles, func(f string) Matcher {
		ignore, err := ignore.CompileIgnoreFile(f)

		if err != nil {
			panic(err)
		}

		return wrapper(*ignore)
	})
}

type wrapper ignore.GitIgnore

func (w wrapper) MatchString(f string) bool {
	var i ignore.GitIgnore = ignore.GitIgnore(w)
	return i.MatchesPath(f)
}
