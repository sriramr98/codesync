package git

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"gitub.com/sriramr98/codesync/database"
	"gitub.com/sriramr98/codesync/object"
)

func (g Git) Commit(author object.Person, committer object.Person, commitMessage string) (string, error) {
	//TODO: Support for nested directories
	files, err := g.workspace.ListFiles(false, false)
	if err != nil {
		return "", err
	}

	writtenFiles, err := g.writeBlobs(files)
	if err != nil {
		return "", err
	}

	treeId, err := g.writeTree(writtenFiles)
	if err != nil {
		return "", err
	}

	parentId, err := g.GetHeadRef()
	if err != nil {
		parentId = ""
	}

	log.Printf("ParentID: %s\n", parentId)

	//TODO: Add support for parent
	commit := object.CommitObject{
		TreeId:        treeId,
		Author:        author,
		ParentId:      parentId,
		Committer:     committer,
		CommitMessage: commitMessage,
	}
	commitObj, err := database.NewObject(commit)
	if err != nil {
		return "", err
	}
	commitId, err := g.Write(commitObj)
	if err != nil {
		return "", nil
	}

	err = g.writeCommitToHEAD(commitId)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d files written to GIT. TreeID: %s, CommitID: %s", len(writtenFiles), treeId, commitId), err
}

func (g Git) writeBlobs(files []os.DirEntry) (map[string]string, error) {
	writtenFiles := make(map[string]string)
	for _, file := range files {
		if !file.IsDir() {
			content, err := os.ReadFile(path.Join(g.workspace.Path(), file.Name()))
			if err != nil {
				return nil, err
			}
			blob := object.BlobObject{
				Content: string(content),
			}

			dbObject, err := database.NewObject(blob)
			if err != nil {
				return nil, err
			}

			id, err := g.Database.Write(dbObject)
			if err != nil {
				return nil, err
			}

			writtenFiles[file.Name()] = id
		}
	}

	return writtenFiles, nil
}

func (g Git) writeTree(files map[string]string) (string, error) {
	var treeEntries []object.TreeNode
	for fileName, fileId := range files {
		treeEntries = append(treeEntries, object.TreeNode{
			FileName: fileName,
			Sha:      fileId,
			//TODO: Extract mode from file
			Mode: "100644",
		})
	}
	dirTree := object.TreeObject{
		Nodes: treeEntries,
	}
	dbObj, err := database.NewObject(dirTree)
	if err != nil {
		return "", err
	}

	return g.Write(dbObj)
}

func (g Git) writeCommitToHEAD(commitId string) error {
	headRefBytes, err := os.ReadFile(path.Join(g.gitPath, "HEAD"))
	if err != nil {
		return err
	}

	// Remove the last \n character
	headRef := strings.Split(string(headRefBytes[:len(headRefBytes)-1]), "ref: ")[1]
	fmt.Printf("Updating Ref %s\n", headRef)
	refPath := path.Join(g.gitPath, headRef)
	return os.WriteFile(refPath, []byte(commitId), 0644)
}
