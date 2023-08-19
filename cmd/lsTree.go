package cmd

import (
	"github.com/spf13/cobra"
	"gitub.com/sriramr98/codesync/repo"
	"log"
)

// lsTreeCmd represents the lsTree command
var lsTreeCmd = &cobra.Command{
	Use:   "ls-tree",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		objectSha := args[0]
		gitRepo := repo.Repo{}

		err := gitRepo.PrintTree(objectSha)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsTreeCmd)
}
