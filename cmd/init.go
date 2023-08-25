package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"gopkg.in/ini.v1"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

var ErrUnableToInitialize = errors.New("unable to initialize repo")
var ErrAlreadyInitialized = errors.New("repo already initialized")

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new git repo",
	Long:  `Initialize a new git repo, ignores if already a git repo`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var path string
		if len(args) > 0 {
			path = args[0]
		} else {
			workingDir, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}
			path = workingDir
		}

		err := initGitRepo(path)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Empty repo initialized")
	},
}

func initGitRepo(rootPath string) error {
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
	gitConfig, err := getConfig()
	if err != nil {
		return ErrUnableToInitialize
	}

	folders := []string{
		path.Join(gitPath, "branches"),
		path.Join(gitPath, "objects"),
		path.Join(gitPath, "refs", "tags"),
		path.Join(gitPath, "refs", "heads"),
	}
	filesToCreate := map[string][]byte{
		path.Join(gitPath, "description"): []byte("Unnamed repository; edit this file 'description' to name the repository."),
		path.Join(gitPath, "HEAD"):        []byte("ref: refs/heads/master\n"),
		path.Join(gitPath, "config"):      gitConfig,
	}

	for _, folder := range folders {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return ErrUnableToInitialize
		}
	}

	for file, content := range filesToCreate {
		if err := os.WriteFile(file, content, 0644); err != nil {
			return ErrUnableToInitialize
		}
	}

	return nil
}

func getConfig() ([]byte, error) {
	config := ini.Empty()
	core, err := config.NewSection("core")
	if err != nil {
		return nil, err
	}

	// the version of the gitdir format. 0 means the initial format, 1 the same with extensions. If > 1, git will panic; rs will only accept 0.
	_, err = core.NewKey("repositoryformatversion", "0")
	if err != nil {
		return nil, err
	}
	// disable tracking of file mode (permissions) changes in the work tree.
	_, err = core.NewKey("filemode", "true")
	if err != nil {
		return nil, err
	}
	// indicates that this repository has a worktree. Git supports an optional worktree key which indicates the location of the worktree, if not ..; rs doesn’t.
	_, err = core.NewKey("bare", "false")
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	writer := bufio.NewWriter(buf)
	if _, err := config.WriteTo(writer); err != nil {
		return nil, err
	}

	err = writer.Flush()
	if err != nil {
		return nil, err
	}

	return io.ReadAll(buf)
}

func init() {

	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
