package workspace

import (
	"log"
	"os"
	"path"
	"strings"
)

type Tree struct {
	RootPath string
	Name     string
	Files    []os.DirEntry
	SubTrees []Tree
}

// BuildTree constructs a nested tree representation of the workspace
func (w Workspace) BuildTree() (Tree, error) {
	return w.buildTree(w.rootPath)
}

func (w Workspace) buildTree(rootPath string) (Tree, error) {
	files, err := os.ReadDir(rootPath)
	if err != nil {
		return Tree{}, err
	}

	folders := strings.Split(rootPath, "/")
	folderName := folders[len(folders)-1]

	root := &Tree{
		Name:     folderName,
		RootPath: rootPath,
	}

	for _, file := range files {
		if file.Name() == ".git" {
			continue
		}
		if file.IsDir() {
			log.Printf("building tree for folder %s\n", file.Name())
			subTree, err := w.buildTree(path.Join(rootPath, file.Name()))
			if err != nil {
				log.Println(err)
				return Tree{}, err
			}
			root.SubTrees = append(root.SubTrees, subTree)
		} else {
			root.Files = append(root.Files, file)
		}
	}

	return *root, nil

}

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
