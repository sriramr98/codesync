package object

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

type TreeNode struct {
	Mode string
	Path string
	Sha  string
}

type TreeObject struct {
	Nodes []TreeNode
}

func NewTreeObject(object GitObject) (TreeObject, error) {
	return treeObjectParser(object.Content)
}

func (b TreeObject) GetType() string {
	return "tree"
}

func (b TreeObject) Encode() (string, error) {
	sort.Slice(b.Nodes, func(i, j int) bool {
		return treeNodeSortKey(b.Nodes[i]) < treeNodeSortKey(b.Nodes[j])
	})

	var result bytes.Buffer
	for _, node := range b.Nodes {
		result.WriteString(node.Mode)
		result.WriteString(" ")
		result.WriteString(node.Path)
		result.WriteByte('\x00')

		sha, err := hex.DecodeString(node.Sha)
		if err != nil {
			return "", err
		}
		err = binary.Write(&result, binary.BigEndian, sha)
		if err != nil {
			return "", err
		}
	}

	return result.String(), nil
}

func treeNodeSortKey(node TreeNode) string {
	if strings.HasPrefix(string(node.Mode), "10") {
		return node.Path
	} else {
		return node.Path + "/"
	}
}

func treeObjectParser(objContent string) (TreeObject, error) {
	start := 0
	contentBytes := []byte(objContent)
	max := len(contentBytes)

	treeObj := &TreeObject{}

	for start < max {
		treeNode, endIndex, err := treeLeafParser(contentBytes[start:])
		if err != nil {
			return TreeObject{}, err
		}
		treeObj.Nodes = append(treeObj.Nodes, treeNode)
		// the end of the previous line will be the start for the new line
		start += endIndex
	}

	return *treeObj, nil
}

func treeLeafParser(data []byte) (TreeNode, int, error) {
	fileModeEndIndex := bytes.IndexByte(data, byte(' '))
	if fileModeEndIndex != 5 && fileModeEndIndex != 6 {
		fmt.Println(string(data))
		return TreeNode{}, -1, errors.New("invalid tree format")
	}

	mode := data[:fileModeEndIndex]
	if len(mode) == 5 {
		// normalize to 6 byte mode
		// ex, if mode was 12345, it becomes <spc>12345 which is a 6 byte mode including space character
		mode = append([]byte(" "), mode...)
	}

	// The path ends with a null character
	pathEndIndex := bytes.IndexByte(data, '\x00')
	path := data[fileModeEndIndex+1 : pathEndIndex]

	// we fetch 20 bytes from the null character considered as SHA and convert to HEX string
	sha := hex.EncodeToString(data[pathEndIndex+1 : pathEndIndex+21])
	// we pad it with 0's upto 40 bytes
	sha = fmt.Sprintf("%040s", sha)

	return TreeNode{Mode: string(mode), Path: string(path), Sha: sha}, pathEndIndex + 21, nil
}
