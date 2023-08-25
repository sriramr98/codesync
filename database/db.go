package database

import (
	"bytes"
	"errors"
	"fmt"
	"gitub.com/sriramr98/codesync/libs/sha"
	"gitub.com/sriramr98/codesync/libs/zlib"
	"io"
	"io/fs"
	"os"
	"path"
	"strconv"
)

type Database[T any] interface {
	Read(string) (T, error)
	Write(T) (string, error)
}

type GitDB struct {
	path string // the path to store objects
}

func NewGitDB(path string) GitDB {
	return GitDB{
		path: path,
	}
}

// ID has two parts -> first two characters are folder name and the rest are filename
func (d GitDB) Read(id string) (Object, error) {
	objectFile, err := findFileFromID(d.path, id)
	if err != nil {
		return Object{}, err
	}
	defer objectFile.Close()

	data, err := zlib.UnCompress(objectFile)
	if err != nil {
		return Object{}, err
	}

	// finding the first space to detect the object type
	typeEndIndex := bytes.Index(data, []byte(" "))
	objectType := data[0:typeEndIndex]

	// finding the null character to obtain object size
	nullEndIndex := bytes.Index(data, []byte("\x00"))
	size, err := strconv.Atoi(string(data[typeEndIndex+1 : nullEndIndex]))
	if err != nil {
		return Object{}, err
	}

	// if parsed size doesn't match the actual size of remaining data, return error
	if size != len(data)-nullEndIndex-1 {
		return Object{}, fmt.Errorf("size mismatch")
	}

	objectContent := string(data[nullEndIndex+1:])
	return Object{Type: string(objectType), Content: objectContent}, nil
}

func (d GitDB) Write(data Object) (string, error) {
	objectContent := data.Encode()

	shaHex := sha.ConvertToShaHex(objectContent)

	objectFolderName := shaHex[0:2]
	objectFileName := shaHex[2:]

	foldePath := path.Join(d.path, objectFolderName)

	err := os.MkdirAll(foldePath, 0755)
	if err != nil {
		return "", err
	}

	var b bytes.Buffer
	zlib.Compress(objectContent, &b)
	if err != nil {
		return "", err
	}

	dataToWrite, err := io.ReadAll(&b)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(
		path.Join(foldePath, objectFileName),
		dataToWrite,
		0644,
	)
	if err != nil {
		return "", err
	}

	return shaHex, nil
}

func findFileFromID(dbPath string, id string) (fs.File, error) {
	objectFolderName := id[0:2]
	objectFile := id[2:]

	objectDirFS := os.DirFS(path.Join(dbPath, objectFolderName))
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
