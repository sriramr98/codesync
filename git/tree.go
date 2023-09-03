package git

import (
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
		treeStr, err := node.ToString()
		if err != nil {
			return err
		}

		fmt.Println(treeStr)
	}

	return nil
}
