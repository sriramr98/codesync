package git

import (
	"errors"
	"gitub.com/sriramr98/codesync/database"
	"gitub.com/sriramr98/codesync/workspace"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Git struct {
	database.Database[database.Object]

	gitPath   string
	workspace workspace.Workspace
}

func NewGit(dirPath string) (Git, error) {
	gitPath, err := findGitDir(dirPath)
	if err != nil {
		return Git{}, err
	}

	projectRootPath, found := strings.CutSuffix(gitPath, "/.git")
	if !found {
		return Git{}, errors.New("invalid GIT dir")
	}
	if err != nil {
		return Git{}, err
	}

	objectsPath := path.Join(gitPath, "objects")
	ws, err := workspace.NewWorkspace(projectRootPath)
	if err != nil {
		return Git{}, err
	}

	return Git{
		Database:  database.NewGitDB(objectsPath),
		gitPath:   gitPath,
		workspace: ws,
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
