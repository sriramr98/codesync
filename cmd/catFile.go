/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/spf13/cobra"
	"gitub.com/sriramr98/codesync/repo"
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

		repo := repo.Repo{}
		projectRootDir, err := repo.FindRootDir(currentDirPath)
		if err != nil {
			log.Fatal(err)
		}

		gitPath := path.Join(projectRootDir, ".git")
		object, err := repo.ReadObject(gitPath, args[0])
		if err != nil {
			log.Fatal(err)
		}

		if prettyPrint {
			fmt.Println(object.Content)
			return
		}

		if typePrint {
			fmt.Println(object.Type)
		}

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
