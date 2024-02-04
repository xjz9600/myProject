package compress

type Compress interface {
	Code() uint8
	Compression(data []byte) ([]byte, error)
	UnCompression(data []byte) ([]byte, error)
}

type NoCompress struct {
}

func (n NoCompress) Code() uint8 {
	return 1
}

func (n NoCompress) Compression(data []byte) ([]byte, error) {
	return data, nil
}

func (n NoCompress) UnCompression(data []byte) ([]byte, error) {
	return data, nil
}
