package object

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"sort"
	"strings"
)

// This represents a single tree structure written to GIT
type TreeNode struct {
	Mode     string
	FileName string
	Sha      string
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

// Add a / to sort key so that directories are sorted after files
func treeNodeSortKey(node TreeNode) string {
	if strings.HasPrefix(node.Mode, "10") {
		return node.FileName
	} else {
		return node.FileName + "/"
	}
}
