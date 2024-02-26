package proxy

import (
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	"time"

	"github.com/anthony-dong/golang/pkg/logs"
)

type Proxy struct {
	ListenAddr Addr
	DialAddr   Addr
	Handler    Handler
}

type Handler interface {
	HandlerConn(readConn net.Conn, dialAddr Addr) error
}

func NewProxy(listen, dial string, h Handler) *Proxy {
	if h == nil {
		panic(fmt.Errorf(`invalid handler`))
	}
	return &Proxy{
		ListenAddr: MustParseAddr(listen),
		DialAddr:   MustParseAddr(dial),
		Handler:    h,
	}
}

func (p *Proxy) Run() error {
	listen, err := p.ListenAddr.Listen()
	if err != nil {
		return err
	}
	retryNum := 0
	for {
		conn, err := listen.Accept()
		if err != nil {
			if retryNum > 3 {
				return err
			}
			retryNum = retryNum + 1
			time.Sleep(time.Millisecond * time.Duration(retryNum))
			continue
		}
		retryNum = 0
		go func(conn net.Conn) {
			defer func() {
				if r := recover(); r != nil {
					logs.Error("conn [%s] find panic: %v, stack:\n%s", conn.RemoteAddr(), r, debug.Stack())
				}
				if err := conn.Close(); err != nil {
					logs.Error("conn [%s] close find err: %v", conn.RemoteAddr(), err)
					return
				}
				logs.Debug("conn [%s] closed", conn.RemoteAddr())
			}()
			if err := p.Handler.HandlerConn(conn, p.DialAddr); err != nil {
				if strings.Contains(err.Error(), "unknown protocol") {
					return
				}
				logs.Error("conn [%s] find err: %v", conn.RemoteAddr(), err)
			}
		}(conn)
	}
}
