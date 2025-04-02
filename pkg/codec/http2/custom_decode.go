package http2

import (
	"encoding/binary"
	"fmt"

	"golang.org/x/net/http2"
)

func DecodeGrpcDataFrame(dataFrame *http2.DataFrame) (_ *GrpcDataFrame, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf(`%v`, r)
		}
	}()
	const dataFrameHeaderLen = 5
	result := GrpcDataFrame{
		Type: "grpc",
	}
	data := dataFrame.Data()
	if len(data) < dataFrameHeaderLen {
		return nil, fmt.Errorf("invalid grpc data frame")
	}
	if data[0] == 1 {
		result.IsCompressed = true
	}
	size := binary.BigEndian.Uint32(data[1:])
	result.Size = size
	result.Data = data[dataFrameHeaderLen:][:size]
	return &result, nil
}
