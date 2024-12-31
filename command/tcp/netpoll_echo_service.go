package tcp

import (
	"context"
	"time"

	"github.com/cloudwego/netpoll"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

func newNetPollEchoService(ctx context.Context, addr string) error {
	parseAddr, err := utils.ParseAddr(addr)
	if err != nil {
		return err
	}
	logs.CtxInfo(ctx, "[netpoll] listener addr: %s", parseAddr)
	listen, err := netpoll.CreateListener(parseAddr.Network, parseAddr.Addr)
	if err != nil {
		return err
	}
	defer listen.Close()
	ops := make([]netpoll.Option, 0)
	eventLoop, err := netpoll.NewEventLoop(
		func(ctx context.Context, conn netpoll.Connection) error {
			handleConn(context.Background(), &wrapperStdConn{conn})
			return nil
		},
		ops...,
	)
	if err != nil {
		return err
	}
	if err := eventLoop.Serve(listen); err != nil {
		return err
	}
	return nil
}

type wrapperStdConn struct {
	netpoll.Connection
}

func (s *wrapperStdConn) SetReadDeadline(t time.Time) error {
	sub := t.Sub(time.Now())
	if sub <= 0 {
		return nil
	}
	return s.SetReadTimeout(sub)
}

func (s *wrapperStdConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (s *wrapperStdConn) SetDeadline(t time.Time) error {
	return nil
}
