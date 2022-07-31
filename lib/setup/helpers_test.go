package setup

import (
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/codethread/dotty/lib/fp"
	"github.com/codethread/dotty/lib/git"
	"github.com/stretchr/testify/assert"
)

var fixtures = map[string]string{
	".gitignore": `.private
.config/exercism`,
	".gitignore_global": `DS_Store`,
}

func TestParseGitIgnores(t *testing.T) {
	testDir := createTestFixtures(t, fixtures)

	ignores := ParseGitIgnores(testDir)

	matches := fp.Filter(
		func(patterns git.Pattern) bool {
			println("I WAS URUOIJ")
			return patterns.Matches(".private")
		}, ignores)

	print("ahhh")
	assert.Len(t, matches, 1)
	// assert.Equal(ignores, })
}

func TestGetAllLinkableFiles(t *testing.T) {
	t.Skip()
	assert.NotPanics(t,
		func() {
			var r []regexp.Regexp

			GetAllLinkableFiles("", &r)
		},
	)
}

func createTestFixtures(t *testing.T, data map[string]string) string {
	dir := t.TempDir()

	for file, data := range data {
		os.WriteFile(path.Join(dir, file), []byte(data), 0777)
	}

	return dir
}
