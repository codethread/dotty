package test

import (
	"os"
	"path"
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

func CreateTestFixtures(t *testing.T, data map[string]string) string {
	dir := t.TempDir()

	for file, data := range data {
		os.WriteFile(path.Join(dir, file), []byte(data), 0777)
	}

	return dir
}
