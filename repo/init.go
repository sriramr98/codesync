package repo

import (
	"errors"
	"log"
	"os"
	"path"
	"path/filepath"
)

var ErrUnableToInitialize = errors.New("unable to initialize repo")
var ErrAlreadyInitialized = errors.New("repo already initialized")

func (r Repo) Init(rootPath string) error {

	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		return err
	}

	gitPath := path.Join(rootPath, ".git")
	gitPath, err := filepath.Abs(gitPath)
	if err != nil {
		return err
	}

	log.Printf("Checking for git path %s\n", gitPath)
	if _, err := os.Stat(gitPath); err == nil {
		// git repo already exists
		return ErrAlreadyInitialized
	}

	log.Println("Initializing a new repo at " + rootPath)

	if err := os.MkdirAll(path.Join(gitPath, "branches"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := os.MkdirAll(path.Join(gitPath, "objects"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := os.MkdirAll(path.Join(gitPath, "refs", "tags"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := os.MkdirAll(path.Join(gitPath, "refs", "heads"), 0755); err != nil {
		return ErrUnableToInitialize
	}

	if err := os.WriteFile(
		path.Join(gitPath, "description"),
		[]byte("Unnamed repository; edit this file 'description' to name the repository."),
		0644,
	); err != nil {
		return ErrUnableToInitialize
	}

	if err := os.WriteFile(
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

	return os.WriteFile(path.Join(gitPath, "config"), []byte(gitConfig), 0644)
}
