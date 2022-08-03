package setup

import (
	"path"

	"github.com/codethread/dotty/lib/git"
)

type FileTree struct {
	Dir   string
	Files []string
	Dirs  []FileTree
}

// getAllLinkableFiles finds all files at DIR, ignoring IGNORED
func GetAllLinkableFiles(dir string, ignored git.Patterns) FileTree {
	panic("not implemented")
}

// func GetIgnoredPatterns(dir string, ignoreFiles []string) git.Patterns {
// 	return fp.Pipe2(
// 		fp.MapFilterErr(func(f string) (git.Patterns, error) {
// 			return git.ParseFile(path.Join(dir, f))
// 		}),

// 		fp.Reduce(git.Patterns{}, git.Patterns.Concat),
// 	)(ignoreFiles)
// }

func GetIgnoredPatterns(dir string, ignoreFiles []string) (p git.Patterns) {
	for _, f := range ignoreFiles {
		ignore, err := git.ParseFile(path.Join(dir, f))
		if err == nil {
			p = p.Concat(ignore)
		}
	}

	return
}
