package repo

import (
	"errors"
	"fmt"
	"path"

	"github.com/spf13/afero"
)

var ErrUnableToInitialize = errors.New("unable to initialize repo")
var ErrAlreadyInitialized = errors.New("repo already initialized")

func (r Repo) Init(rootPath string) error {

	gitPath := path.Join(rootPath, ".git")

	if _, err := r.fs.Stat(gitPath); err == nil {
		return ErrAlreadyInitialized
	}

	fmt.Println("Initializing a new repo at " + rootPath)

	if err := r.fs.MkdirAll(path.Join(gitPath, "branches"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := r.fs.MkdirAll(path.Join(gitPath, "objects"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := r.fs.MkdirAll(path.Join(gitPath, "refs", "tags"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := r.fs.MkdirAll(path.Join(gitPath, "refs", "heads"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := afero.WriteFile(
		r.fs,
		path.Join(gitPath, "description"),
		[]byte("Unnamed repository; edit this file 'description' to name the repository."),
		0644,
	); err != nil {
		return ErrUnableToInitialize
	}

	if err := afero.WriteFile(
		r.fs,
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

	return afero.WriteFile(r.fs, path.Join(gitPath, "config"), []byte(gitConfig), 0644)
}
