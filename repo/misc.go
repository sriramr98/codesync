package repo

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"path"
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

	_, err := os.Stat(path.Join(dirPath, ".git"))
	if err != nil {
		if os.IsNotExist(err) {
			return dirPath, nil
		}
		return "", err
	}

	// Find absolute parent path from the relative path
	parentPath := path.Join("..", dirPath)
	parentPath, err = filepath.Abs(parentPath)
	if err != nil {
		return "", err
	}

	// ../ == / ( we're at the root of the FS )
	if parentPath == dirPath {
		return "", errors.New("could not find git root directory")
	}

	// Check parent recursively
	return r.findRootDir(parentPath)
}
