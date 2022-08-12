package lib

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/codethread/dotty/lib/fp"
	ignore "github.com/sabhiram/go-gitignore"
)

var naps int

func Setup(config SetupConfig) {
	fmt.Print("\n\n\n\n")
	start := time.Now()
	// teardown
	ignored := getIgnoredPatterns(config.gitignores)
	allIgnored := append(ignored, config.ignored...)
	fmt.Println("start")

	c := make(chan FileTree)
	go getAllLinkableFiles(config.from, allIgnored, c)
	files := <-c

	asStr := ToGOB64(files)
	os.WriteFile(config.historyFile, []byte(asStr), 0644)
	// files.Walk(Visitor{
	// 	file: func(dir string, file string) { fmt.Println("->", dir, file) },
	// 	dir:  func(dir string) { fmt.Println("dd", dir) },
	// })

	// for each file, link it into home
	// for all success, add them to a teardown file
	// for all failures, present a warning in the console

	// fmt.Println(files)
	fmt.Println("duration", time.Since(start).Milliseconds())

}

// go binary encoder
func ToGOB64(m FileTree) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

// go binary decoder
func FromGOB64(str string) FileTree {
	m := FileTree{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(`failed base64 Decode`, err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return m
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
	ignore3, err := regexp.Compile("^node_modules$")
	if err != nil {
		panic(err)
	}

	matcher = append(matcher, ignore, ignore2, ignore3)

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

func (ft FileTree) Debug() string {
	return p(ft, 0, 2)
}

type Visitor struct {
	File func(dir string, file string)
	Dir  func(string)
}

func (ft FileTree) Walk(vistor Visitor) {
	vistor.Dir(ft.Dir)

	for _, f := range ft.Files {
		vistor.File(ft.Dir, f)
	}

	for _, fts := range ft.Dirs {
		fts.Walk(vistor)
	}
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
