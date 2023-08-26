package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitub.com/sriramr98/codesync/git"
	"log"
)

// TODO: Add show-ref for a specific reference
var showRefCmd = &cobra.Command{
	Use:   "show-ref",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		gitRepo, err := git.NewRepo(".")
		if err != nil {
			log.Fatal(err)
		}

		refs, err := gitRepo.FetchRefs()
		if err != nil {
			log.Fatal(err)
		}

		for _, ref := range refs {
			fmt.Printf("%s %s\n", ref.SHA, ref.Path)
		}
	},
}

func init() {
	rootCmd.AddCommand(showRefCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showRefCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showRefCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
