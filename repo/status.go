package repo

import (
	"errors"
	"fmt"
	"os"
	"path"
)

func (r Repo) Status() (string, error) {

	pwd, err := os.Getwd()
	fmt.Printf("Current Path %s\n", pwd)
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(path.Join(pwd, ".git")); err != nil {
		if os.IsNotExist(err) {
			rootDir, err := r.findRootDir(pwd)
			if err != nil {
				return "", errors.New("not a git repository")
			}
			return "Git Root at " + rootDir, nil
		}
		return "", err
	}

	return "Git Root at " + pwd, nil
}
