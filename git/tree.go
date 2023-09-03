package git

import (
	"bytes"
	"errors"
	"fmt"

	"gitub.com/sriramr98/codesync/parsers"
)

func (g Git) PrintTree(objectSha string) error {
	gitObj, err := g.Read(objectSha)
	if err != nil {
		return err
	}

	treeObj, err := parsers.ParseTree([]byte(gitObj.Content))
	if err != nil {
		return err
	}

	//fmt.Println(treeObj)
	for _, node := range treeObj.Nodes {

		objMode := node.Mode[0:2]
		objType := byteToType([]byte(objMode))

		if objType == "" {
			return errors.New("invalid tree node")
		}

		fmt.Printf("%s %s %s\t%s\n", node.Mode, objType, node.Sha, node.FileName)
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
