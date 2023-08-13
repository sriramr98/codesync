package zlib

import (
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

func Compress(data []byte, w io.Writer) {
	zio := zlib.NewWriter(w)
	zio.Write(data)
	zio.Close()
}
