package lib

import (
	"github.com/codethread/dotty/lib/fp"
	ignore "github.com/sabhiram/go-gitignore"
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
	return fp.PromiseAll(ignoreFiles, func(f string) Matcher {
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
