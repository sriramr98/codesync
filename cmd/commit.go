package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitub.com/sriramr98/codesync/git"
	"gitub.com/sriramr98/codesync/object"
	"log"
	"os"
	"time"
)

var authorName string
var authorEmail string

var commitMessage string

// commitCmd represents the commit command
var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("commit called")
		gitRepo, err := git.NewGit(".")
		if err != nil {
			log.Fatal(err)
		}

		if commitMessage == "" {
			log.Fatal("Commit Message Required")
		}
		if authorName == "" {
			name, hasName := os.LookupEnv("GIT_AUTHOR_NAME")
			if !hasName {
				log.Fatal("Author Name required")
			}
			authorName = name
		}
		if authorEmail == "" {
			email, hasEmail := os.LookupEnv("GIT_AUTHOR_EMAIL")
			if !hasEmail {
				log.Fatal("Author Email required")
			}
			authorEmail = email
		}

		now := time.Now()
		_, offset := now.Zone()
		timestamp := fmt.Sprintf("%d %d", now.Unix(), offset)
		author := object.Person{
			Name:      authorName,
			Email:     authorEmail,
			Timestamp: timestamp,
		}

		// We don't support custom committer
		result, err := gitRepo.Commit(author, author, commitMessage)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(result)
	},
}

func init() {

	commitCmd.PersistentFlags().StringVarP(&authorName, "authorName", "n", "", "Author Name for the commit")
	commitCmd.PersistentFlags().StringVarP(&authorEmail, "authorEmail", "e", "", "Author Email for the commit")
	commitCmd.PersistentFlags().StringVarP(&commitMessage, "message", "m", "", "Commit Message")

	rootCmd.AddCommand(commitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
