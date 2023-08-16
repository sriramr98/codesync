package object

type TagObject struct {
	Content string
}

func NewTagObject(object GitObject) TagObject {
	return TagObject{object.Content}
}

func (b TagObject) GetContent() string {
	return b.Content
}

func (b TagObject) GetType() string {
	return "blob"
}
