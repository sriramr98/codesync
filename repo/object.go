package repo

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strconv"

	zlib "gitub.com/sriramr98/codesync/utils"
)

type GitObject struct {
	Type    string
	Content string
}

func (r Repo) ReadObject(gitFolderPath string, sha string) (GitObject, error) {
	objectFile, err := r.findFileFromSHA(gitFolderPath, sha)
	if err != nil {
		return GitObject{}, err
	}
	defer objectFile.Close()

	data, err := zlib.UnCompress(objectFile)
	if err != nil {
		return GitObject{}, err
	}

	// finding the first space to detect the object type
	typeEndIndex := bytes.Index(data, []byte(" "))
	objectType := data[0:typeEndIndex]

	// finding the null character to obtain object size
	nullEndIndex := bytes.Index(data, []byte("\x00"))
	size, err := strconv.Atoi(string(data[typeEndIndex+1 : nullEndIndex]))
	if err != nil {
		return GitObject{}, err
	}

	// if parsed size doesn't match the actual size of remaining data, return error
	if size != len(data)-nullEndIndex-1 {
		return GitObject{}, fmt.Errorf("size mismatch")
	}

	objectContent := string(data[nullEndIndex+1:])

	return GitObject{Type: string(objectType), Content: objectContent}, nil

}

func (r Repo) findFileFromSHA(rootFolderPath string, sha string) (fs.File, error) {
	objectFolderName := sha[0:2]
	objectFile := sha[2:]

	objectDirFS := os.DirFS(path.Join(rootFolderPath, "objects", objectFolderName))
	files, err := fs.Glob(objectDirFS, fmt.Sprintf("%s*", objectFile))
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, errors.New("invalid SHA, file not found")
	}

	if len(files) > 1 {
		return nil, errors.New("invalid SHA, too many files found for pattern")
	}

	return objectDirFS.Open(files[0])
}
