package tcp

import (
	"context"
	"io"
	"net"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/logs"
)

func NewEchoServiceCommand() (*cobra.Command, error) {
	addr := ""
	command := &cobra.Command{
		Use:   "echo_service",
		Short: "create a tcp echo service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return newEchoService(cmd.Context(), addr)
		},
	}
	command.Flags().StringVar(&addr, "addr", ":8080", "监听地址")
	return command, nil
}

func newEchoService(ctx context.Context, addr string) error {
	parseAddr, err := utils.ParseAddr(addr)
	if err != nil {
		return err
	}
	logs.CtxInfo(ctx, "listener addr: %s", parseAddr)
	listen, err := parseAddr.Listen()
	if err != nil {
		return err
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}
		go func() {
			logs.CtxInfo(ctx, "receive conn %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
			defer func() {
				_ = conn.Close()
				logs.CtxInfo(ctx, "close conn %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
			}()
			if err := echoHandler(ctx, conn); err != nil {
				if err == io.EOF {
					return
				}
				logs.CtxError(ctx, "%s conn find err: %v", conn.RemoteAddr(), err)
			}
		}()
	}
}

func echoHandler(ctx context.Context, conn net.Conn) error {
	buffer := make([]byte, 1024)
	for {
		readSize, err := conn.Read(buffer)
		if err != nil {
			return err
		}
		if _, err := conn.Write(buffer[:readSize]); err != nil {
			return err
		}
	}
}
