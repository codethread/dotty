package setup

import (
	"path"
	"regexp"

	"github.com/codethread/dotty/lib/git"
)

type FileTree struct {
	Dir   string
	Files []string
	Dirs  []FileTree
}

// getAllLinkableFiles finds all files at DIR, ignoring IGNORED
func GetAllLinkableFiles(dir string, ignored *[]regexp.Regexp) FileTree {
	panic("not implemented")
}

func ParseGitIgnores(dir string) []git.Pattern {
	ignore, err := git.ParseGitIgnoreFile(path.Join(dir, ".gitignore"))

	if err != nil {
		panic(err)
	}

	global, err := git.ParseGitIgnoreFile(path.Join(dir, ".gitignore_global"))

	if err != nil {
		panic(err)
	}

	return append(ignore, global...)
}
