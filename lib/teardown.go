package lib

import (
	"fmt"
	"io"
	"os"
)

func Teardown(config SetupConfig) {
	dryRun := config.DryRun
	if dryRun {
		fmt.Println("dry run, no files will be deleted...")
	}

	files := ParseFromFile[FileTree](config.HistoryFile)
	// delete all files
	files.Walk(FileTreeVisitor{
		Dir: func(string) {},
		File: func(dir string, file string) {
			target := fromFromPathToTarget(&config, dir, file)
			if dryRun {
				fmt.Println("delete ->", target)
			} else {
				deleteTarget(target)
			}
		},
	})

	if dryRun {
		return
	}

	// directories are walked first, so after deleting all files
	// we check for empty dirs
	files.Walk(FileTreeVisitor{
		File: func(string, string) {},
		Dir: func(dir string) {
			target := fromFromPathToTarget(&config, dir)
			empty, _ := isEmpty(target)
			if empty {
				deleteTarget(target)
			}
		},
	})
}

func deleteTarget(target string) {
	err := os.Remove(target)
	if err != nil {
		fmt.Println("could not remove", target, err)
	}
}

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
