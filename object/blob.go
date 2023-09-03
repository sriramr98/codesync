package object

import "fmt"

type BlobObject struct {
	Content string
}

func (b BlobObject) GitType() string {
	return "blob"
}

func (b BlobObject) Encode() (string, error) {
	return b.Content, nil
}

func (b BlobObject) Print() {
	fmt.Println(b.Content)
}
