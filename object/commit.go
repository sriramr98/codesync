package object

import (
	"fmt"
)

type CommitObject struct {
	TreeId        string
	ParentId      string
	Author        Person
	Committer     Person
	CommitMessage string
}

type Person struct {
	Name      string
	Email     string
	Timestamp string
}

// Format Formats the object the way git recognises it
func (p Person) Format() string {
	return fmt.Sprintf("%s <%s> %s", p.Name, p.Email, p.Timestamp)
}

func (b CommitObject) GitType() string {
	return "commit"
}

// Encode Converts the content of the CommitObject as a git recognised format to compress and create git object
func (b CommitObject) Encode() (string, error) {
	authorStr := b.Author.Format()
	committerStr := b.Committer.Format()
	//output := fmt.Sprintf("tree %s\nparent %s\nauthor %s\ncommitter %s\n\n%s", b.TreeId, b.ParentId, authorStr, committerStr, b.CommitMessage)

	output := fmt.Sprintf("tree %s", b.TreeId)
	if b.ParentId != "" {
		output = fmt.Sprintf("%s\nparent %s", output, b.ParentId)
	}
	output = fmt.Sprintf("%s\nauthor %s\ncommitter%s\n\n%s", output, authorStr, committerStr, b.CommitMessage)

	return output, nil
}
