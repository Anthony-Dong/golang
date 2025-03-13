package compress

type Type uint8

const (
	NopType    Type = 0
	GzipType   Type = 1
	SnappyType Type = 2
)

type Compressor interface {
	Decompress(data []byte) ([]byte, error)
	Compress(data []byte) ([]byte, error)
}

var _ Compressor = Nop{}

type Nop struct{}

func (n Nop) Decompress(data []byte) ([]byte, error) {
	return data, nil
}

func (n Nop) Compress(data []byte) ([]byte, error) {
	return data, nil
}

func NewCompressor(t Type) Compressor {
	switch t {
	case GzipType:
		return Gzip{}
	case SnappyType:
		return Snappy{}
	}
	return Nop{}
}
