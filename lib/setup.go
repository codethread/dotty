package lib

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/codethread/dotty/lib/fp"
	ignore "github.com/sabhiram/go-gitignore"
)

func Setup(config SetupConfig) {
	fmt.Print("\n\n\n\n")
	// teardown
	ignored := getIgnoredPatterns(config.gitignores)
	allIgnored := append(ignored, config.ignored...)
	files := getAllLinkableFiles(config.from, allIgnored)
	// for each file, link it into home
	// for all success, add them to a teardown file
	// for all failures, present a warning in the console

	fmt.Println(files)

}

type SetupConfig struct {
	dryRun      bool
	from        string
	to          string
	ignored     []Matcher
	gitignores  []string
	historyFile string
}

func expand(home string) func(string) string {
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

func BuildSetupConfig(flags Flags, implicitConfig ImplicitConfig) SetupConfig {
	e := expand(implicitConfig.Home)
	var matcher []Matcher

	ignore, err := regexp.Compile("^_.*")
	ignore2, err := regexp.Compile("^.git$")
	if err != nil {
		panic(err)
	}

	matcher = append(matcher, ignore, ignore2)

	return SetupConfig{
		dryRun:      true,
		from:        e("~/PersonalConfigs"),
		to:          e("~/test"),
		gitignores:  []string{e("~/PersonalConfigs/.gitignore_global"), e("~/PersonalConfigs/.gitignore")},
		ignored:     matcher,
		historyFile: e("~/.dotty"),
	}
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

func getAllLinkableFiles(dir string, ignores Matchers) (ft FileTree) {
	fS := os.DirFS(dir)
	dirE, err := fs.ReadDir(fS, ".")

	if err != nil {
		panic(err)
	}

	ft.Dir = dir

	for _, f := range dirE {
		n := f.Name()
		if ignores.Matches(n) {
			continue
		} else if f.IsDir() {
			ft.Dirs = append(ft.Dirs, getAllLinkableFiles(path.Join(dir, n), ignores))
		} else {
			ft.Files = append(ft.Files, n)
		}
	}

	return

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

type Matcher interface {
	MatchString(f string) bool
}

type FileTree struct {
	Dir   string
	Files []string
	Dirs  []FileTree
}

func (ft FileTree) String() string {
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
