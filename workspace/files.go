package workspace

import (
	"os"
	"path"
)

// ListFiles Lists all files and folders in the workspace.
// If recursive is false, only lists files and folders in root dir
// If recursive is true, includeDirs is ignored ( will be true in this case )
func (w Workspace) ListFiles(recursive bool, includeDirs bool) ([]os.DirEntry, error) {
	if recursive {
		includeDirs = true
	}
	files, err := listFiles(w.rootPath)
	if err != nil {
		return nil, err
	}

	if !recursive {
		return files, nil
	}

	for _, file := range files {
		if file.IsDir() && includeDirs {
			dirFiles, err := listFiles(path.Join(w.rootPath, file.Name()))
			if err != nil {
				return nil, err
			}
			files = append(files, dirFiles...)
		}

		files = append(files, file)
	}

	return files, nil
}

func listFiles(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}
