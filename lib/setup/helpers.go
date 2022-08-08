package setup

import (
	"path"

	"github.com/sabhiram/go-gitignore"
)

type Matcher interface {
	MatchesPath(f string) bool
}

type FileTree struct {
	Dir   string
	Files []string
	Dirs  []FileTree
}

// getAllLinkableFiles finds all files at DIR, ignoring IGNORED
func GetAllLinkableFiles(dir string, ignores []Matcher) FileTree {

	panic("not implemented")
}

func GetIgnoredPatterns(dir string, ignoreFiles []string) (ignores []Matcher) {
	for _, f := range ignoreFiles {
		ignore, err := ignore.CompileIgnoreFile(path.Join(dir, f))
		if err != nil {
			ignores = append(ignores, ignore)
		}
	}
	return
}
