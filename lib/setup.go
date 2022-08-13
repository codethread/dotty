package lib

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
)

var naps int

func Setup(config SetupConfig) {
	dryRun := config.DryRun

	if dryRun {
		fmt.Println("dry run, no files will be changed...")
	}

	Teardown(config)

	ignored := GetIgnoredPatterns(config.gitignores)
	allIgnored := append(ignored, config.ignored...)

	c := make(chan FileTree)
	go getAllLinkableFiles(config.From, allIgnored, c)
	files := <-c

	files = files.Filter(func(dir string, file string) bool {
		target := fromFromPathToTarget(&config, dir, file)

		if fileAlreadyExists(target) {
			fmt.Fprintf(os.Stderr, "file %s already exists, you'll need to remove it manually\n", target)
			return false
		}
		return true
	})

	if !dryRun {
		StringifyToFile(files, config.HistoryFile)
	}

	// for each file, link it into home
	files.Walk(FileTreeVisitor{
		File: func(dir string, file string) {
			from := path.Join(dir, file)
			target := fromFromPathToTarget(&config, dir, file)
			if dryRun {
				fmt.Println("symlink ->", target)
			} else {
				err := os.Symlink(from, target)
				if err != nil {
					fmt.Println("could not symlink", from, "=>", target)
				}
			}
		},
		Dir: func(dir string) {
			target := fromFromPathToTarget(&config, dir)
			if dryRun {
				fmt.Println("create =>", target)
			} else {
				os.MkdirAll(target, os.ModePerm)
			}
		},
	})
}

func getAllLinkableFiles(dir string, ignores Matchers, c chan FileTree) {
	ft := FileTree{
		Dir: dir,
	}
	fS := os.DirFS(dir)

	dirE, err := fs.ReadDir(fS, ".")

	if err != nil {
		panic(err)
	}

	var dirs []string

	for _, f := range dirE {
		name := f.Name()
		if ignores.Matches(name) {
			continue
		} else if f.IsDir() {
			dirs = append(dirs, name)
		} else {
			ft.Files = append(ft.Files, name)
		}
	}

	c2 := make(chan FileTree, len(dirs))

	for _, name := range dirs {
		go getAllLinkableFiles(path.Join(dir, name), ignores, c2)
	}

	for range dirs {
		ft.Dirs = append(ft.Dirs, <-c2)
	}

	c <- ft
}

func fileAlreadyExists(file string) bool {
	fi, err := os.Lstat(file)

	// no file
	if err != nil {
		return false
	}

	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		// symlink
		return false
	}

	// file exists
	return true
}

func fromFromPathToTarget(config *SetupConfig, paths ...string) string {
	paths[0] = strings.Replace(paths[0], config.From, config.To, 1)
	return path.Join(paths...)
}
