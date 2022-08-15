package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func TestFiles(config SetupConfig, files []string) {
	var invalidFiles []string
	var validFiles []string

	dotfiles, invalidFiles := NormaliseTestFiles(config, files)

	ignores := GetAllIgnoredPatterns(config)

	for _, dotfile := range dotfiles {
		if fileDoesntExist(dotfile, config.From) {
			invalidFiles = append(invalidFiles, dotfile)
		} else if isFileDotty(ignores, dotfile) {
			validFiles = append(validFiles, dotfile)
		} else {
			invalidFiles = append(invalidFiles, dotfile)
		}
	}

	for _, dotfile := range invalidFiles {
		fmt.Fprintf(os.Stderr, "invalid %s\n", dotfile)
	}

	if len(invalidFiles) != 0 {
		for _, dotfile := range validFiles {
			fmt.Println("valid", dotfile)
		}
		os.Exit(1)
	}
}

func fileDoesntExist(file string, from string) bool {
	_, err := os.Lstat(filepath.Join(from, file))

	if err != nil {
		// file doesnt exist
		return true
	}

	return false
}

func isFileDotty(ignores Matchers, file string) bool {
	parts := strings.Split(file, "/")
	for _, part := range parts {
		if ignores.Matches(part) {
			return false
		}
	}
	return true
}

func NormaliseTestFiles(config SetupConfig, files []string) (validFiles []string, invalidFiles []string) {
	e := ExpandHome(config.ImplicitConfig.Home)

	for _, file := range files {
		expanded := e(file)
		if !strings.HasPrefix(expanded, "/") {
			validFiles = append(validFiles, expanded)
		} else if strings.HasPrefix(expanded, config.From) {
			prefix := len(config.From) + 1
			validFiles = append(validFiles, expanded[prefix:])
		} else {
			invalidFiles = append(invalidFiles, file)
		}
	}

	return
}
