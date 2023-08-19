package repo

import (
	"bytes"
	"errors"
	"fmt"
	"gitub.com/sriramr98/codesync/repo/object"
)

func (r Repo) PrintTree(objectSha string) error {
	gitPath, err := r.FindGitDir(".")
	if err != nil {
		return err
	}

	gitObj, err := r.ReadObject(gitPath, objectSha)
	if err != nil {
		return err
	}

	treeObj, err := object.NewTreeObject(gitObj)
	if err != nil {
		return err
	}

	for _, node := range treeObj.Nodes {

		objMode := node.Mode[0:2]
		objType := byteToType([]byte(objMode))

		if objType == "" {
			return errors.New("invalid tree node")
		}

		fmt.Printf("%s %s %s\t%s\n", node.Mode, objType, node.Sha, node.Path)
	}

	return nil
}

func byteToType(data []byte) string {
	if bytes.Equal(data, []byte("04")) {
		return "tree"
	} else if bytes.Equal(data, []byte("10")) {
		return "blob" // this is a normal blob file
	} else if bytes.Equal(data, []byte("12")) {
		return "blob" // this is a symlink pointing to a blob
	} else if bytes.Equal(data, []byte("16")) {
		return "commit"
	} else {
		return ""
	}
}
