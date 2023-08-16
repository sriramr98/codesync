package repo

import (
	"errors"
	"fmt"
	"os"
)

func (r Repo) Status() (string, error) {

	pwd, err := os.Getwd()
	fmt.Printf("Current Path %s\n", pwd)
	if err != nil {
		return "", err
	}

	rootDir, err := r.FindGitDir(pwd)
	if err != nil {
		return "", errors.New("not a git repository")
	}
	return "Git Root at " + rootDir, nil

}
