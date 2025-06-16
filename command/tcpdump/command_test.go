package tcpdump

import (
	"context"
	"testing"

	"github.com/anthony-dong/golang/pkg/tcpdump"
)

func TestDecodeHTTP(t *testing.T) {
	file := "test/out_kitex_big.pcap"

	option := tcpdump.NewDecodeOptions(tcpdump.WithDecoder(tcpdump.NewHttp2Decoder()), tcpdump.WithMsgWriter(tcpdump.NewConsoleLogMessageWriter([]tcpdump.MessageType{
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

	option := tcpdump.NewDecodeOptions(tcpdump.WithDecoder(tcpdump.NewThriftDecoder()), tcpdump.WithMsgWriter(tcpdump.NewConsoleLogMessageWriter([]tcpdump.MessageType{
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
