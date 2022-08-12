package lib

import (
	"fmt"
	"path"
	"strings"
)

type FileTree struct {
	Dir   string
	Files []string
	Dirs  []FileTree
}

type FileTreeVisitor struct {
	File func(dir string, file string)
	Dir  func(string)
}

func (ft FileTree) Walk(vistor FileTreeVisitor) {
	vistor.Dir(ft.Dir)

	for _, f := range ft.Files {
		vistor.File(ft.Dir, f)
	}

	for _, fts := range ft.Dirs {
		fts.Walk(vistor)
	}
}

func (ft FileTree) Filter(predicate func(dir string, file string) bool) FileTree {
	var newFiles []string
	for _, f := range ft.Files {
		if predicate(ft.Dir, f) {
			newFiles = append(newFiles, f)
		}
	}
	ft.Files = newFiles

	for i, fts := range ft.Dirs {
		ft.Dirs[i] = fts.Filter(predicate)
	}

	return ft
}

func (ft FileTree) Debug() string {
	return p(ft, 0, 2)
}

func p(ft FileTree, indent int, indentation int) string {
	i := strings.Repeat(" ", indent*indentation)

	out := fmt.Sprintf("\n%s%s/", i, path.Base(ft.Dir))

	for _, f := range ft.Files {
		out = fmt.Sprintf("%s\n%s%s", out, i, f)
	}

	for _, f := range ft.Dirs {
		out = fmt.Sprintf("%s%s", out, p(f, indent+1, indentation))
	}
	return out
}
