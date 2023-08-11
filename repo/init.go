package repo

import (
	"errors"
	"fmt"
	"path"
)

var ErrUnableToInitialize = errors.New("unable to initialize repo")

func (r Repo) Init(rootPath string) error {

	gitPath := path.Join(rootPath, ".git")

	if err := r.fs.Exists(gitPath); err == nil {
		fmt.Println("Already initialized")
		return nil
	}

	fmt.Println("Initializing a new repo at " + rootPath)

	if err := r.fs.Mkdir(path.Join(gitPath, "branches"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := r.fs.Mkdir(path.Join(gitPath, "objects"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := r.fs.Mkdir(path.Join(gitPath, "refs", "tags"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := r.fs.Mkdir(path.Join(gitPath, "refs", "heads"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := r.fs.WriteFile(
		path.Join(gitPath, "description"),
		[]byte("Unnamed repository; edit this file 'description' to name the repository."),
		0644,
	); err != nil {
		return ErrUnableToInitialize
	}

	if err := r.fs.WriteFile(
		path.Join(gitPath, "HEAD"),
		[]byte("ref: refs/heads/master\n"),
		0644,
	); err != nil {
		return ErrUnableToInitialize
	}

	gitConfig, err := r.getConfig()
	if err != nil {
		return ErrUnableToInitialize
	}

	return r.fs.WriteFile(path.Join(gitPath, "config"), []byte(gitConfig), 0644)
}
