package setup

import (
	"path"

	"github.com/codethread/dotty/lib/fp"
	. "github.com/codethread/dotty/lib/git"
)

type FileTree struct {
	Dir   string
	Files []string
	Dirs  []FileTree
}

// getAllLinkableFiles finds all files at DIR, ignoring IGNORED
func GetAllLinkableFiles(dir string, ignored Patterns) FileTree {
	panic("not implemented")
}

func GetIgnoredPatterns(dir string, ignoreFiles []string) Patterns {
	patterns := fp.FilterMap(func(f string) (Patterns, bool) {
		ignore, err := ParseGitIgnoreFile(path.Join(dir, f))
		return ignore, err == nil

	}, ignoreFiles)

	return fp.Reduce(Patterns{}, patterns, func(acc Patterns, cur Patterns) Patterns {
		return acc.Concat(cur)
	})
}

func ParseGitIgnores(dir string) Patterns {
	ignore, err := ParseGitIgnoreFile(path.Join(dir, ".gitignore"))

	if err != nil {
		panic(err)
	}

	global, err := ParseGitIgnoreFile(path.Join(dir, ".gitignore_global"))

	if err != nil {
		panic(err)
	}

	return ignore.Concat(global)
}
