package object

import (
	"bytes"
	"strconv"
)

type ObjectSerializer interface {
	Serialize() []byte
}

type GitObject struct {
	Type    string
	Content string
}

// object format
// <type> <content_len>\x00<content>
// \x00 is null character
func (b GitObject) Encode() []byte {
	bytesBuffer := new(bytes.Buffer)
	bytesBuffer.Write([]byte(b.Type))
	bytesBuffer.Write([]byte(" "))
	bytesBuffer.Write([]byte(strconv.Itoa(len(b.Content))))
	bytesBuffer.Write([]byte("\x00"))
	bytesBuffer.Write([]byte(b.Content))

	return bytesBuffer.Bytes()
}
