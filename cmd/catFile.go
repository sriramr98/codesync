package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gitub.com/sriramr98/codesync/git"
	"gitub.com/sriramr98/codesync/object"
	"gitub.com/sriramr98/codesync/parsers"
)

var prettyPrint bool = false
var typePrint bool = false

// catFileCmd represents the catFile command
var catFileCmd = &cobra.Command{
	Use:   "cat-file",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		if !prettyPrint && !typePrint {
			fmt.Println("Please specify either -p or -t")
			return
		}

		if prettyPrint && typePrint {
			fmt.Println("Please specify either -p or -t")
			return
		}

		currentDirPath, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}

		repo, err := git.NewGit(currentDirPath)
		if err != nil {
			log.Fatal(err)
		}

		dbObj, err := repo.Read(args[0])
		if err != nil {
			log.Fatal(err)
		}

		if typePrint {
			fmt.Println(dbObj.Type)
			return
		}

		var gitObj object.GitObject
		switch dbObj.Type {
		case "tree":
			gitObj, err = parsers.ParseTree([]byte(dbObj.Content))
		case "commit":
			gitObj, err = parsers.ParseCommit([]byte(dbObj.Content))
		case "blob":
			gitObj = object.BlobObject{
				Content: dbObj.Content,
			}
		}

		if err != nil {
			log.Fatal(err)
		}

		gitObj.Print()
	},
}

func init() {

	catFileCmd.PersistentFlags().BoolVarP(&prettyPrint, "pretty", "p", false, "pretty print the output")
	catFileCmd.PersistentFlags().BoolVarP(&typePrint, "type", "t", false, "print the object type")

	rootCmd.AddCommand(catFileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// catFileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// catFileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
