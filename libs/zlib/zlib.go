package zlib

import (
	"bytes"
	"compress/zlib"
	"io"
)

func UnCompress(r io.Reader) ([]byte, error) {
	zio, err := zlib.NewReader(r)
	if err != nil {
		return nil, err
	}
	defer zio.Close()
	return io.ReadAll(zio)
}

func Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	zio := zlib.NewWriter(&b)
	defer zio.Close()

	zio.Write(data)
	return b.Bytes(), nil
}
