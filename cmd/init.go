package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"gitub.com/sriramr98/codesync/repo"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new git repo",
	Long:  `Initialize a new git repo, ignores if already a git repo`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := repo.Repo{}

		path := "."
		if len(args) > 0 {
			path = args[0]
		}

		err := repo.Init(path)

		if err != nil {
			log.Fatal(err)
		}
	},
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
