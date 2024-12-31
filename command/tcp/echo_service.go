package tcp

import (
	"context"
	"errors"
	"io"
	"net"

	"github.com/cloudwego/netpoll"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/logs"
)

func NewEchoServiceCommand() (*cobra.Command, error) {
	addr := ""
	isNetpoll := false
	command := &cobra.Command{
		Use:   "echo_service",
		Short: "create a tcp echo service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if isNetpoll {
				return newNetPollEchoService(cmd.Context(), addr)
			}
			return newEchoService(cmd.Context(), addr)
		},
	}
	command.Flags().StringVar(&addr, "addr", ":8080", "监听地址")
	command.Flags().BoolVar(&isNetpoll, "netpoll", false, "use github.com/cloudwego/netpoll")
	return command, nil
}

func newEchoService(ctx context.Context, addr string) error {
	parseAddr, err := utils.ParseAddr(addr)
	if err != nil {
		return err
	}
	logs.CtxInfo(ctx, "[std] listener addr: %s", parseAddr)
	listen, err := parseAddr.Listen()
	if err != nil {
		return err
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			return err
		}
		go handleConn(ctx, conn)
	}
}

func handleConn(ctx context.Context, conn net.Conn) {
	logs.CtxInfo(ctx, "receive conn %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
	defer func() {
		_ = conn.Close()
		logs.CtxInfo(ctx, "close conn %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
	}()
	if err := echoHandler(ctx, conn); err != nil {
		if err == io.EOF || errors.Is(err, netpoll.ErrEOF) {
			return
		}
		logs.CtxError(ctx, "%s conn find err: %v", conn.RemoteAddr(), err)
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
