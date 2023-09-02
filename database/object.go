package database

import (
	"bytes"
	"gitub.com/sriramr98/codesync/object"
	"strconv"
)

type Object struct {
	Type    string
	Content string
}

func NewObject(gitObj object.GitObject) (Object, error) {
	content, err := gitObj.Encode()
	if err != nil {
		return Object{}, err
	}
	return Object{
		Type:    gitObj.GitType(),
		Content: content,
	}, nil
}

// Encode object format
// <type> <content_len>\x00<content>
// \x00 is null character
func (o Object) Encode() []byte {
	bytesBuffer := new(bytes.Buffer)
	bytesBuffer.Write([]byte(o.Type))
	bytesBuffer.Write([]byte(" "))
	bytesBuffer.Write([]byte(strconv.Itoa(len(o.Content))))
	bytesBuffer.Write([]byte("\x00"))
	bytesBuffer.Write([]byte(o.Content))

	return bytesBuffer.Bytes()
}
