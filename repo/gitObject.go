package repo

import (
	"bytes"
	"strconv"
)

type GitObject struct {
	Type    string
	Content string
}

// object format
// <type> <content_len>\x00<content>
// \x00 is null character
func (gitObj GitObject) Encode() []byte {
	bytesBuffer := new(bytes.Buffer)
	bytesBuffer.Write([]byte(gitObj.Type))
	bytesBuffer.Write([]byte(" "))
	bytesBuffer.Write([]byte(strconv.Itoa(len(gitObj.Content))))
	bytesBuffer.Write([]byte("\x00"))
	bytesBuffer.Write([]byte(gitObj.Content))

	return bytesBuffer.Bytes()
}
