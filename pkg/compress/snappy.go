package compress

import "github.com/golang/snappy"

var _ Compressor = Snappy{}

type Snappy struct {
}

func (s Snappy) Decompress(data []byte) ([]byte, error) {
	return snappy.Decode(nil, data)
}

func (s Snappy) Compress(data []byte) ([]byte, error) {
	rest := snappy.Encode(nil, data)
	return rest, nil
}
