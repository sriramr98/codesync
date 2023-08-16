package object

type BlobObject struct {
	content string
}

func NewBlobObject(gitObject GitObject) BlobObject {
	return BlobObject{
		gitObject.Content,
	}
}

func (b BlobObject) GetContent() string {
	return b.content
}

func (b BlobObject) GetType() string {
	return "blob"
}
