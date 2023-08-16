package repo

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"strconv"

	"gitub.com/sriramr98/codesync/libs/sha"
	"gitub.com/sriramr98/codesync/libs/zlib"
	"gitub.com/sriramr98/codesync/repo/object"
)

func (r Repo) ReadObject(gitFolderPath string, sha string) (object.GitObject, error) {
	objectFile, err := r.findFileFromSHA(gitFolderPath, sha)
	if err != nil {
		return object.GitObject{}, err
	}
	defer objectFile.Close()

	data, err := zlib.UnCompress(objectFile)
	if err != nil {
		fmt.Println("Uncompress err")
		return object.GitObject{}, err
	}

	// finding the first space to detect the object type
	typeEndIndex := bytes.Index(data, []byte(" "))
	objectType := data[0:typeEndIndex]

	// finding the null character to obtain object size
	nullEndIndex := bytes.Index(data, []byte("\x00"))
	size, err := strconv.Atoi(string(data[typeEndIndex+1 : nullEndIndex]))
	if err != nil {
		return object.GitObject{}, err
	}

	// if parsed size doesn't match the actual size of remaining data, return error
	if size != len(data)-nullEndIndex-1 {
		return object.GitObject{}, fmt.Errorf("size mismatch")
	}

	objectContent := string(data[nullEndIndex+1:])
	return object.GitObject{Type: string(objectType), Content: objectContent}, nil
}

func (r Repo) WriteObject(gitFolderPath string, data object.GitObject) (string, error) {
	objectContent := data.Encode()

	shaHex := sha.ConvertToShaHex(objectContent)

	objectFolderName := shaHex[0:2]
	objectFileName := shaHex[2:]

	foldePath := path.Join(gitFolderPath, "objects", objectFolderName)

	err := os.MkdirAll(foldePath, 0755)
	if err != nil {
		return "", err
	}

	// file, err := os.Open(path.Join(foldePath, objectFileName))
	// if err != nil {
	// 	return "", err
	// }

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
