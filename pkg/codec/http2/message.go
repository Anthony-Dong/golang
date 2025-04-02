package http2

import (
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

type GrpcDataFrame struct {
	Type         string `json:"type"`
	IsCompressed bool   `json:"is_compressed"`
	Size         uint32 `json:"size"`
	Data         []byte `json:"data"`
}

type HTTP2Frame struct {
	Type          string              `json:"type,omitempty"`
	Flags         http2.Flags         `json:"flags,omitempty"`
	Length        uint32              `json:"length,omitempty"`
	StreamID      uint32              `json:"stream_id,omitempty"`
	Data          interface{}         `json:"data,omitempty"`
	Headers       []hpack.HeaderField `json:"headers,omitempty"`
	Settings      map[string]uint32   `json:"settings,omitempty"`
	TransportInfo *TransportInfo      `json:"transport_info,omitempty"`
}

type TransportInfo struct {
	FromAddr string `json:"from_addr,omitempty"`
	ToAddr   string `json:"to_addr,omitempty"`
}

/**
客户端发起的流：客户端发起的新流会被分配奇数编号的流 ID，初始值通常从 1 开始，后续新发起的流依次递增，例如 1、3、5 等。当客户端向服务器发送一个新的 gRPC 请求时，就会创建一个新的流并分配一个奇数的流 ID。
服务器发起的流：服务器发起的新流会被分配偶数编号的流 ID，起始值一般为 2，后续依次递增，如 2、4、6 等。这种情况相对较少，例如在服务器主动向客户端推送消息时会使用偶数 ID 的流。
*/

func ConvertToFrame(frame http2.Frame) *HTTP2Frame {
	header := frame.Header()
	result := HTTP2Frame{
		Type:     header.Type.String(),
		Flags:    header.Flags,
		Length:   header.Length,
		StreamID: header.StreamID,
	}
	switch v := frame.(type) {
	case *http2.DataFrame:
		dataFrame, err := DecodeGrpcDataFrame(v)
		if err == nil {
			result.Data = dataFrame
		} else {
			result.Data = v.Data()
		}
	case *http2.MetaHeadersFrame:
		result.Headers = v.Fields
	case *http2.SettingsFrame:
		result.Settings = FormatSettings(v)
	}
	return &result
}

func FormatSettings(s *http2.SettingsFrame) map[string]uint32 {
	result := map[string]uint32{}
	_ = s.ForeachSetting(func(setting http2.Setting) error {
		result[setting.ID.String()] = uint32(setting.Val)
		return nil
	})
	return result
}
