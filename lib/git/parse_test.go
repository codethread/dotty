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

		testDir := test.CreateTestFixtures(t, fixtures)

		patterns, _ := ParseGitIgnoreFile(path.Join(testDir, ".gitignore"))

		shouldMatch := []string{
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
		}

		shouldNotMatch := []string{
			"#any",
			"any#",
			"flychecksomething.el",
			".ideanested",
			".config/bar",
			".config/execism",
			"private",
			"cheese",
		}

		for _, tc := range shouldMatch {
			Convey(fmt.Sprintf("Pattern %v should be matched", tc), func() {
				So(patterns.Matches(tc), ShouldBeTrue)
			})
		}

		for _, tc := range shouldNotMatch {
			Convey(fmt.Sprintf("Pattern %v should NOT be matched", tc), func() {
				So(patterns.Matches(tc), ShouldBeFalse)
			})
		}
	})

}
