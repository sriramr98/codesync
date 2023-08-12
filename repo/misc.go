package repo

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

func (r *Repo) getConfig() ([]byte, error) {
	config := ini.Empty()
	core, err := config.NewSection("core")
	if err != nil {
		return nil, err
	}

	// the version of the gitdir format. 0 means the initial format, 1 the same with extensions. If > 1, git will panic; rs will only accept 0.
	core.NewKey("repositoryformatversion", "0")
	// disable tracking of file mode (permissions) changes in the work tree.
	core.NewKey("filemode", "true")
	// indicates that this repository has a worktree. Git supports an optional worktree key which indicates the location of the worktree, if not ..; rs doesn’t.
	core.NewKey("bare", "false")

	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	if _, err := config.WriteTo(writer); err != nil {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	return io.ReadAll(buf)
}

func (r *Repo) findRootDir(dirPath string) (string, error) {
	if !filepath.IsAbs(dirPath) {
		path, err := filepath.Abs(dirPath)
		if err != nil {
			return "", err
		}
		dirPath = path
	}

	// Check if dirPath has a folder .git
	if _, err := os.Stat(filepath.Join(dirPath, ".git")); err == nil {
		return dirPath, nil
	}

	parentPath := filepath.Dir(dirPath)
	if parentPath == dirPath {
		return "", errors.New("could not find .git directory")
	}

	// Check parent recursively
	return r.findRootDir(parentPath)
}
