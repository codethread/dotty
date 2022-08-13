package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/codethread/dotty/test"
	. "github.com/smartystreets/goconvey/convey"
)

var gitignores = map[string]string{
	"from/.gitignore":        ".private\n.config/exercism",
	"from/.gitignore_global": "bar",
}

var fromFiles = map[string]string{
	"from/.private": "",
	"from/bar/a":    "",
	"from/a":        "",
	"from/b":        "",
	"from/c/a":      "",
	"from/c/b":      "",
	"from/c/bar":    "",
	"from/d/a/a":    "",
	"from/d/a/b":    "",
	"from/d/bar/a":  "",
}

var shouldExist = map[string]string{
	"to/a":     "",
	"to/b":     "",
	"to/c/a":   "",
	"to/c/b":   "",
	"to/d/a/a": "",
	"to/d/a/b": "",
}

var shouldNotExist = map[string]string{
	"from/.private": "",
	"from/bar/a":    "",
	"from/d/bar/a":  "",
}

func TestSetup(t *testing.T) {
	Convey("Setup", t, func() {
		testDir := test.CreateFixtures(t, gitignores, fromFiles)
		e := joinBase(testDir)

		ignore, err := regexp.Compile("^.gitignore$")
		if err != nil {
			panic(err)
		}

		config := SetupConfig{
			gitignores: []string{
				e("from/.gitignore"),
				e("from/.gitignore_global"),
			},
			ignored:     []Matcher{ignore},
			HistoryFile: e(".dotty"),
			From:        e("from"),
			To:          e("to"),
		}

		Setup(config)

		for file := range shouldExist {
			target := filepath.Join(testDir, file)
			Convey(fmt.Sprintf("file '%s' should exist", file), func() {
				_, err := os.Stat(target)
				So(err, ShouldBeNil)
			})
		}

		for file := range shouldNotExist {
			target := filepath.Join(testDir, file)
			Convey(fmt.Sprintf("file '%s' should NOT exist", file), func() {
				_, err := os.Stat(target)
				So(err, ShouldBeNil)
			})
		}
	})
}

func joinBase(dir string) func(string) string {
	return func(path string) string {
		return filepath.Join(dir, path)
	}
}
