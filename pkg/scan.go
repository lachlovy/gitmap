package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

var skipDirs = map[string]bool{
	"node_modules": true,
}

func shouldSkipDir(folder string) bool {
	return skipDirs[folder]
}

func findGitRepositoryInDirectory(directory string, allRepositories []string) []string {
	allDirs, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return allRepositories
	}

	subDirNames := make([]string, 0)

	for _, dir := range allDirs {
		if dir.IsDir() && dir.Name() == ".git" {
			return append(allRepositories, directory)
		}
	}

	for _, dir := range allDirs {
		if dir.IsDir() && !shouldSkipDir(dir.Name()) {
			subDir := filepath.Join(directory, dir.Name())
			allRepositories = findGitRepositoryInDirectory(subDir, allRepositories)
		}
	}

	for _, subDirName := range subDirNames {
		subDir := filepath.Join(directory, subDirName)
		allRepositories = findGitRepositoryInDirectory(subDir, allRepositories)
	}
	return allRepositories
}

func ScanGitRepositories(directory string) []string {
	var allRepositories = make([]string, 0)
	allRepositories = findGitRepositoryInDirectory(directory, allRepositories)
	return allRepositories
}
