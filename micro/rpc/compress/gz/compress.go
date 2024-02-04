package gz

import (
	"bytes"
	"compress/gzip"
	"io"
)

type GzipCompress struct {
}

func (g GzipCompress) Code() uint8 {
	return 2
}

func (g GzipCompress) Compression(data []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write(data); err != nil {
		return nil, err
	}
	if err := gz.Flush(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (g GzipCompress) UnCompression(data []byte) ([]byte, error) {
	res := new(bytes.Buffer)
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		var out []byte
		return out, err
	}
	defer reader.Close()
	_, err = io.Copy(res, reader)
	return res.Bytes(), err
}
