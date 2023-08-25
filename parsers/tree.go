package parsers

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"gitub.com/sriramr98/codesync/object"
)

func ParseTree(data []byte) (object.TreeObject, error) {
	start := 0
	max := len(data)

	treeObj := &object.TreeObject{}

	for start < max {
		treeNode, endIndex, err := treeLeafParser(data[start:])
		if err != nil {
			return object.TreeObject{}, err
		}
		treeObj.Nodes = append(treeObj.Nodes, treeNode)
		// the end of the previous line will be the start for the new line
		start += endIndex
	}

	return *treeObj, nil
}

func treeLeafParser(data []byte) (object.TreeNode, int, error) {
	fileModeEndIndex := bytes.IndexByte(data, byte(' '))
	if fileModeEndIndex != 5 && fileModeEndIndex != 6 {
		fmt.Println(string(data))
		return object.TreeNode{}, -1, errors.New("invalid tree format")
	}

	mode := data[:fileModeEndIndex]
	if len(mode) == 5 {
		// normalize to 6 byte mode
		// ex, if mode was 12345, it becomes <spc>12345 which is a 6 byte mode including space character
		mode = append([]byte("0"), mode...)
	}

	// The path ends with a null character
	pathEndIndex := bytes.IndexByte(data, '\x00')
	path := data[fileModeEndIndex+1 : pathEndIndex]

	// we fetch 20 bytes from the null character considered as SHA and convert to HEX string
	sha := hex.EncodeToString(data[pathEndIndex+1 : pathEndIndex+21])

	return object.TreeNode{Mode: string(mode), Path: string(path), Sha: sha}, pathEndIndex + 21, nil
}
