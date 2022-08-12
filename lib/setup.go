package lib

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/codethread/dotty/lib/fp"
	ignore "github.com/sabhiram/go-gitignore"
)

var naps int

func Setup(config SetupConfig) {
	fmt.Print("\n\n\n\n")
	start := time.Now()
	// TODO: teardown
	ignored := getIgnoredPatterns(config.gitignores)
	allIgnored := append(ignored, config.ignored...)
	fmt.Println("start")

	c := make(chan FileTree)
	go getAllLinkableFiles(config.From, allIgnored, c)
	files := <-c

	files = files.Filter(func(dir string, file string) bool {
		target := strings.Replace(dir, config.From, config.To, 1)
		targetFile := path.Join(target, file)
		fileExists := fileAlreadyExists(targetFile)
		if fileExists {
			fmt.Println("file", targetFile, "already exists, you'll need to remove it manually")
		}
		return !fileExists
	})

	StringifyToFile(files, config.HistoryFile)

	// for each file, link it into home
	// files.Walk(Visitor{
	// 	file: func(dir string, file string) { fmt.Println("->", dir, file) },
	// 	dir:  func(dir string) { fmt.Println("dd", dir) },
	// })

	// for all success, add them to a teardown file
	// for all failures, present a warning in the console

	fmt.Println("duration", time.Since(start).Milliseconds())

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
