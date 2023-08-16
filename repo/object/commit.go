package object

import (
	"fmt"
	"strings"
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

func NewCommitObject(object GitObject) CommitObject {
	return parseGitObject(object)
}

func (b CommitObject) GetType() string {
	return "blob"
}

// Encode Converts the content of the CommitObject as a git recognised format to compress and create git object
func (b CommitObject) Encode() string {
	authorStr := b.Author.Format()
	committerStr := b.Committer.Format()
	//output := fmt.Sprintf("tree %s\nparent %s\nauthor %s\ncommitter %s\n\n%s", b.TreeId, b.ParentId, authorStr, committerStr, b.CommitMessage)

	output := fmt.Sprintf("tree %s", b.TreeId)
	if b.ParentId != "" {
		output = fmt.Sprintf("%s\nparent %s", output, b.ParentId)
	}
	output = fmt.Sprintf("%s\nauthor %s\ncommitter%s\n\n%s", output, authorStr, committerStr, b.CommitMessage)

	return output
}

// Does not support parsing signed git objects for now
func parseGitObject(object GitObject) CommitObject {
	result := &CommitObject{}

	lines := strings.Split(object.Content, "\n")

	for _, line := range lines {
		lineSplit := strings.SplitN(line, " ", 2)
		if len(line) == 0 {
			continue
		}
		key := lineSplit[0]
		value := lineSplit[1]

		switch key {
		case "tree":
			result.TreeId = value
		case "parent":
			result.ParentId = value
		case "author":
			person := parsePersonDetails(value)
			result.Author = person
		case "committer":
			person := parsePersonDetails(value)
			result.Committer = person
		default:
			if result.CommitMessage == "" {
				result.CommitMessage += line
			} else {
				// since we parse line by line, we remove the \n present in the commit message, so we manually add it
				result.CommitMessage += "\n" + line
			}
		}

	}

	return *result
}

func parsePersonDetails(value string) Person {
	parts := strings.Split(value, "<")
	name := strings.TrimSpace(parts[0])

	secondaryParts := strings.Split(parts[1], ">")
	email := strings.TrimSpace(secondaryParts[0])
	//TODO: Verify the meaning of the component and parse accordingly
	timestamp := strings.TrimSpace(secondaryParts[1])

	return Person{
		Name:      name,
		Email:     email,
		Timestamp: timestamp,
	}
}
