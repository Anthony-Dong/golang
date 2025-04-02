package tcpdump

import (
	"context"
	"testing"

	"github.com/anthony-dong/golang/pkg/tcpdump"
)

func TestDecodeHTTP(t *testing.T) {
	file := "test/out_kitex_big.pcap"

	option := NewDecodeOptions(WithDecoder(tcpdump.NewHttp2Decoder()), WithMsgWriter(tcpdump.NewConsoleLogMessageWriter([]tcpdump.MessageType{
		tcpdump.MessageType_Log,
		tcpdump.MessageType_HTTP2,
	})))
	source, err := NewPacketSource(file, option)
	if err != nil {
		t.Fatal(err)
	}

	if err := DecodePacketSource(context.Background(), source, option); err != nil {
		t.Fatal(err)
	}
}

func TestDecodeThrift(t *testing.T) {
	return
	file := "test/out.pcap"

	option := NewDecodeOptions(WithDecoder(tcpdump.NewThriftDecoder()), WithMsgWriter(tcpdump.NewConsoleLogMessageWriter([]tcpdump.MessageType{
		tcpdump.MessageType_Log,
		tcpdump.MessageType_Thrift,
	})))
	source, err := NewPacketSource(file, option)
	if err != nil {
		t.Fatal(err)
	}
	if err := DecodePacketSource(context.Background(), source, option); err != nil {
		t.Fatal(err)
	}
}
