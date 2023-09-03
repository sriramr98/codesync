package git

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"gitub.com/sriramr98/codesync/database"
	"gitub.com/sriramr98/codesync/object"
	"gitub.com/sriramr98/codesync/workspace"
)

func (g Git) Commit(author object.Person, committer object.Person, commitMessage string) (string, error) {
	//TODO: Support for nested directories
	workspaceTree, err := g.workspace.BuildTree()
	if err != nil {
		return "", err
	}

	log.Printf("Workspace Tree to be written %s", workspaceTree)
	treeId, err := g.recursivelyWriteWorkspace(workspaceTree)

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

	return fmt.Sprintf("TreeID: %s, CommitID: %s", treeId, commitId), err
}

func (g Git) writeBlobs(files []os.DirEntry, filesRootPath string) ([]object.TreeNode, error) {
	writtenFiles := []object.TreeNode{}
	for _, file := range files {
		if !file.IsDir() {
			content, err := os.ReadFile(path.Join(filesRootPath, file.Name()))
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

			fileNode := object.TreeNode{
				Mode:     "100644",
				FileName: file.Name(),
				Sha:      id,
			}
			writtenFiles = append(writtenFiles, fileNode)
		}
	}

	return writtenFiles, nil
}

func (g Git) writeTree(files []object.TreeNode) (string, error) {
	dirTree := object.TreeObject{
		Nodes: files,
	}
	dbObj, err := database.NewObject(dirTree)
	if err != nil {
		return "", err
	}

	return g.Write(dbObj)
}

// writes the entire tree to GIT and returns the sha of the root tree
func (g Git) recursivelyWriteWorkspace(tree workspace.Tree) (string, error) {
	filesWritten, err := g.writeBlobs(tree.Files, tree.RootPath)
	if err != nil {
		return "", err
	}

	for _, folder := range tree.SubTrees {
		log.Printf("Writing subtree for %s with filecount %d", folder.Name, len(folder.Files))
		treeId, err := g.recursivelyWriteWorkspace(folder)
		if err != nil {
			return "", err
		}

		treeNode := object.TreeNode{
			FileName: folder.Name,
			Sha:      treeId,
			Mode:     "40000",
		}
		filesWritten = append(filesWritten, treeNode)
	}

	return g.writeTree(filesWritten)
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
