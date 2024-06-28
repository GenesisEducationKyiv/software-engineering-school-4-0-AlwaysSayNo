package util

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProjectRootPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	rootPath := cwd
	for {
		if isProjectRoot(rootPath) {
			break
		}
		parentDir := filepath.Dir(rootPath)
		if parentDir == rootPath {
			return "", fmt.Errorf("could not find project root")
		}
		rootPath = parentDir
	}

	return rootPath, nil
}

func isProjectRoot(path string) bool {
	_, err := os.Stat(filepath.Join(path, "go.mod"))
	return err == nil
}
