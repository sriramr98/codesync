package object

type TreeObject struct {
	Content string
}

func NewTreeObject(object GitObject) TreeObject {
	return TreeObject{object.Content}
}

func (b TreeObject) GetContent() string {
	return b.Content
}

func (b TreeObject) GetType() string {
	return "blob"
}
