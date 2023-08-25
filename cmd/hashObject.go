package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"gitub.com/sriramr98/codesync/database"
	"gitub.com/sriramr98/codesync/object"
	"log"
	"os"

	"github.com/spf13/cobra"
	"gitub.com/sriramr98/codesync/libs/sha"
	"gitub.com/sriramr98/codesync/repo"
)

var writeObject bool
var readFromStdin bool

func readFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	return string(content), err
}

func readStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	return reader.ReadString('\n')
}

func extractDataToWrite(args []string) (string, error) {
	if readFromStdin {
		return readStdin()
	} else {
		return readFile(args[0])
	}
}

// hashObjectCmd represents the hashObject command
var hashObjectCmd = &cobra.Command{
	Use:   "hash-object",
	Short: "Hashes content",
	Long:  `Converts the input file or content into a Git Hash and optionally writes it to the database`,
	Args:  cobra.MaximumNArgs(1),
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !readFromStdin && len(args) == 0 {
			// not reading data either from stdin or args
			return errors.New("either pass data through --stdin or pass a file path in args")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Flags; stdin %t, write %t\n", readFromStdin, writeObject)
		dataToWrite, err := extractDataToWrite(args)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Data To write %s\n", dataToWrite)

		gitObject := object.BlobObject{
			Content: dataToWrite,
		}

		if !writeObject || readFromStdin {
			// Only calculate HASH and print
			encodedData, err := gitObject.Encode()
			if err != nil {
				log.Fatal(err)
			}
			hash := sha.ConvertToShaHex([]byte(encodedData))
			fmt.Print(hash)
		} else {
			// hash and write
			currentWd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			repo, err := repo.NewRepo(currentWd)
			if err != nil {
				log.Fatal(err)
			}

			hash, err := repo.Write(database.Object{
				Type:    "blob",
				Content: dataToWrite,
			})

			if err != nil {
				log.Fatal(err)
			}
			fmt.Print(hash)
		}

	},
}

func init() {

	hashObjectCmd.PersistentFlags().BoolVarP(&writeObject, "write", "w", false, "write object")
	hashObjectCmd.PersistentFlags().BoolVarP(&readFromStdin, "stdin", "s", false, "read from stdin")

	rootCmd.AddCommand(hashObjectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// hashObjectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// hashObjectCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
