package database

import (
	"bytes"
	"strconv"
)

type Object struct {
	Type    string
	Content string
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
