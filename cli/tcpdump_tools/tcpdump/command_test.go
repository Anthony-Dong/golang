package tcpdump

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/anthony-dong/golang/pkg/tcpdump"
	"github.com/anthony-dong/golang/pkg/utils"
)

func readFile(file string) string {
	dir := utils.GetGoProjectDir()
	return filepath.Join(dir, "command/tcpdump/test", file)
}

func Test_DecodeTCPDump(t *testing.T) {
	ctx := context.Background()
	cfg := tcpdump.NewDefaultConfig()
	t.Run("thrift", func(t *testing.T) {
		if err := Run(ctx, readFile("thrift.pcap"), cfg, DefaultDecoders); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("http", func(t *testing.T) {
		if err := Run(context.Background(), readFile("http1.1.pcap"), cfg, DefaultDecoders); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("http chunked", func(t *testing.T) {
		if err := Run(context.Background(), readFile("http1.1_chunked.pcap"), cfg, DefaultDecoders); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("stick http", func(t *testing.T) {
		if err := Run(context.Background(), readFile("stick_http1.1.pcap"), cfg, DefaultDecoders); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("thrift_ttheader", func(t *testing.T) {
		// thrift_ttheader
		if err := Run(context.Background(), readFile("thrift_ttheader.pcap"), cfg, DefaultDecoders); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("stick_thrift_ttheader", func(t *testing.T) {
		// thrift_ttheader
		if err := Run(ctx, readFile("stick_thrift_ttheader.pcap"), cfg, DefaultDecoders); err != nil {
			t.Fatal(err)
		}
	})
}
