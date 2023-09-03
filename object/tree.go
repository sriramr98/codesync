package object

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
)

// This represents a single tree structure written to GIT
type TreeNode struct {
	Mode     string
	FileName string
	Sha      string
}

func (node TreeNode) ToString() (string, error) {
	objMode := node.Mode[0:2]
	objType := byteToType([]byte(objMode))

	if objType == "" {
		return "", errors.New("invalid tree node")
	}

	return fmt.Sprintf("%s %s %s\t%s", node.Mode, objType, node.Sha, node.FileName), nil
}

type TreeObject struct {
	Nodes []TreeNode
}

func (b TreeObject) GitType() string {
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
		result.WriteString(node.FileName)
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

func (b TreeObject) Print() {
	for _, node := range b.Nodes {
		nodeStr, err := node.ToString()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(nodeStr)
	}
}

// Add a / to sort key so that directories are sorted after files
func treeNodeSortKey(node TreeNode) string {
	if strings.HasPrefix(node.Mode, "10") {
		return node.FileName
	} else {
		return node.FileName + "/"
	}
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
