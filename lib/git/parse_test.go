package git

import (
	"fmt"
	"path"
	"testing"

	"github.com/codethread/dotty/test"
	. "github.com/smartystreets/goconvey/convey"
)

var fixtures = map[string]string{
	".gitignore": `
# files
.private
.config/exercism

# directory
.idea/

# nested
Session.vim
**/Session.vim

# file type
*.rel

# file suffix
flycheck_*.el

# fancy
\#*\#
.\#*

# complex
[._]*.sw[a-p]
`,
}

func TestParseGitIgnoreFile(t *testing.T) {
	Convey("TestParseGitIgnoreFile converts a gitignore file to a valid Patterns struct", t, func() {
		Convey("converts valid gitignore files to Patterns struct", func() {
			testDir := test.CreateFixtures(t, fixtures)

			patterns, err := ParseFile(path.Join(testDir, ".gitignore"))

			Convey("it does not return an error", func() {
				So(err, ShouldEqual, nil)
			})

			for _, tc := range []string{
				".private",
				".config/exercism",
				".idea/nested",
				".idea/deep/nested",
				"some/nested/Session.vim",
				"filetype.rel",
				"nested/filetype.rel",
				"flycheck_something.el",
				"#any#",
				".#any",
				".foo.swp",
			} {
				Convey(fmt.Sprintf("Pattern %v should be matched", tc), func() {
					So(patterns.Matches(tc), ShouldBeTrue)
				})
			}

			for _, tc := range []string{
				"#any",
				"any#",
				"flychecksomething.el",
				".ideanested",
				".config/bar",
				".config/execism",
				"private",
				"cheese",
			} {
				Convey(fmt.Sprintf("Pattern %v should NOT be matched", tc), func() {
					So(patterns.Matches(tc), ShouldBeFalse)
				})
			}

		})

		Convey("ParseFile returns an error when the file is not valid", func() {
			_, err := ParseFile(path.Join("~", "nope"))
			So(err, ShouldNotEqual, nil)
		})
	})

}
