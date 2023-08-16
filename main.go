/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"gitub.com/sriramr98/codesync/repo"
	"gitub.com/sriramr98/codesync/repo/object"
	"log"
)

func main() {
	//cmd.Execute()

	gitRepo := repo.Repo{}
	gitPath, err := gitRepo.FindGitDir("/Users/sriram-r/Desktop/personal/testdir")
	if err != nil {
		log.Fatal(err)
	}

	gitObject, err := gitRepo.ReadObject(gitPath, "26f5")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gitObject.Content)
	commitObj := object.NewCommitObject(gitObject)

	fmt.Println(commitObj.Encode())

}
