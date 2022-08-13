package test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func GetFixtureBase() string {
	ex, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(ex)

	return filepath.Join(exPath, "fixtures")
}

func CreateFixtures(t *testing.T, fixtureList ...map[string]string) string {
	dir := t.TempDir()

	for _, fixtures := range fixtureList {
		for file, contents := range fixtures {
			path := filepath.Join(dir, file)
			os.MkdirAll(filepath.Dir(path), fs.ModePerm)
			os.WriteFile(path, []byte(contents), 0777)
		}
	}

	return dir
}
