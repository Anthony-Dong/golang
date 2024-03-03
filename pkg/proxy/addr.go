package proxy

import (
	"fmt"
	"net"
	"strings"
)

type SimpleAddr struct {
	Network string
	Addr    string
}

type Addr interface {
	Listen() (net.Listener, error)
	Dial() (net.Conn, error)
}

func (a *SimpleAddr) Listen() (net.Listener, error) {
	switch a.Network {
	case "tcp", "udp":
		return net.Listen(a.Network, a.Addr)
	case "unix":
		//addr := &net.UnixAddr{"/dev/shm/shmipc.sock", "unix"}
		addr := &net.UnixAddr{Name: a.Addr, Net: "unix"}
		return net.ListenUnix(a.Network, addr)
	default:
		return nil, fmt.Errorf(`not support netwrok: %s`, a.Network)
	}
}

func (a *SimpleAddr) Dial() (net.Conn, error) {
	switch a.Network {
	case "tcp", "udp":
		return net.Dial(a.Network, a.Addr)
	case "unix":
		addr := &net.UnixAddr{Name: a.Addr, Net: "unix"}
		return net.DialUnix(a.Network, nil, addr)
	default:
		return nil, fmt.Errorf(`not support netwrok: %s`, a.Network)
	}
}

func MustParseAddr(addr string) *SimpleAddr {
	ret := &SimpleAddr{Network: "tcp"}
	splits := strings.SplitN(addr, "://", 2)
	if len(splits) == 1 {
		if strings.Contains(addr, ".socket") || strings.Contains(addr, ".sock") {
			ret.Network = "unix"
		}
		if addr == "" {
			panic(fmt.Errorf(`invalid addr: %s`, addr))
		}
		ret.Addr = addr
		return ret
	}
	ret.Network = splits[0]
	ret.Addr = splits[1]
	if ret.Network == "" || ret.Addr == "" {
		panic(fmt.Errorf(`invalid addr: %s`, addr))
	}
	return ret
}
