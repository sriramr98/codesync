package repo

import (
	"errors"
	"gitub.com/sriramr98/codesync/database"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Repo struct {
	database.Database[database.Object]

	gitPath string
}

func NewRepo(dirPath string) (Repo, error) {
	gitPath, err := findGitDir(dirPath)
	if err != nil {
		return Repo{}, err
	}

	objectsPath := path.Join(gitPath, "objects")

	return Repo{
		Database: database.NewGitDB(objectsPath),
		gitPath:  gitPath,
	}, nil
}

func findGitDir(dirPath string) (string, error) {
	if !filepath.IsAbs(dirPath) {
		absPath, err := filepath.Abs(dirPath)
		if err != nil {
			return "", err
		}
		dirPath = absPath
	}

	if strings.HasSuffix(dirPath, ".git") || strings.HasSuffix(dirPath, ".git/") {
		// given path is git path
		return dirPath, nil
	}

	// Check if dirPath has a folder .git
	gitDir := path.Join(dirPath, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		return gitDir, nil
	}

	parentPath := filepath.Dir(dirPath)
	if parentPath == dirPath {
		return "", errors.New("could not find .git directory")
	}

	// Check parent recursively
	return findGitDir(parentPath)
}
