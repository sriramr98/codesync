package parsers

import (
	"gitub.com/sriramr98/codesync/object"
	"strings"
)

// ParseCommit parses a slice of bytes to a commit object
func ParseCommit(data []byte) (object.GitObject, error) {
	// Does not support parsing signed git objects for now
	result := &object.CommitObject{}

	lines := strings.Split(string(data), "\n")

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

	return *result, nil
}

func parsePersonDetails(value string) object.Person {
	parts := strings.Split(value, "<")
	name := strings.TrimSpace(parts[0])

	secondaryParts := strings.Split(parts[1], ">")
	email := strings.TrimSpace(secondaryParts[0])
	//TODO: Verify the meaning of the component and parse accordingly
	timestamp := strings.TrimSpace(secondaryParts[1])

	return object.Person{
		Name:      name,
		Email:     email,
		Timestamp: timestamp,
	}
}
