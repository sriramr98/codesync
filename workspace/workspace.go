package workspace

import (
	"os"
	"path"
)

type Workspace struct {
	rootPath string
}

func NewWorkspace(rootPath string) (Workspace, error) {
	// if path doesn't have a .git folder, not a valid workspace path
	folder, err := os.Stat(path.Join(rootPath, ".git"))
	if err != nil || !folder.IsDir() {
		return Workspace{}, err
	}

	return Workspace{rootPath: rootPath}, nil
}

func (w Workspace) Path() string {
	return w.rootPath
}
