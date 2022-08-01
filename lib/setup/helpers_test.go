package setup

import (
	"testing"

	"github.com/codethread/dotty/test"
	. "github.com/smartystreets/goconvey/convey"
)

var fixtures = map[string]string{
	".gitignore":        ".private\n.config/exercism",
	".gitignore_global": "bar",
}

func TestGetIgnoredPatterns(t *testing.T) {
	Convey("GetIgnoredPatterns", t, func() {
		testDir := test.CreateTestFixtures(t, fixtures)

		Convey("When no ignore paths are passed", func() {
			patterns := GetIgnoredPatterns(testDir, []string{})
			So(patterns.Patterns, ShouldBeEmpty)
		})

		Convey("When some ignore paths are passed", func() {
			patterns := GetIgnoredPatterns(testDir, []string{".gitignore", ".gitignore_global"})
			So(len(patterns.Patterns), ShouldEqual, 3)
		})

		Convey("When some ignore paths are passed and don't exist", func() {
			patterns := GetIgnoredPatterns(testDir, []string{".gitignore", ".gitignore_global", "missing_file"})
			So(len(patterns.Patterns), ShouldEqual, 3)
		})
	})

}
