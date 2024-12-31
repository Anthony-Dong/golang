package utils

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type SimpleAddr struct {
	Network string
	Addr    string
}

func (a *SimpleAddr) String() string {
	return fmt.Sprintf("%s://%s", a.Network, a.Addr)
}

func (a *SimpleAddr) GetIPPort() (net.IP, int, error) {
	split := strings.SplitN(a.Addr, ":", 2)
	if len(split) != 2 {
		return nil, 0, fmt.Errorf("invalid addr: %s", a.Addr)
	}
	port, err := strconv.ParseInt(split[1], 10, 64)
	if err != nil {
		return nil, 0, fmt.Errorf("invalid addr: %s", a.Addr)
	}
	ip := net.ParseIP(split[0])
	if ip == nil {
		return nil, 0, fmt.Errorf("invalid addr: %s", a.Addr)
	}
	return ip, int(port), nil
}

type Addr interface {
	Listen() (net.Listener, error)
	Dial() (net.Conn, error)
}

type UDPAddr interface {
	ListenUDP() (*net.UDPConn, error)
	DailUDP() (*net.UDPConn, error)
}

var allowedNetwork = map[string]bool{
	"tcp": true,
	//"udp":  true,
	"unix": true,
}

func (a *SimpleAddr) ListenUDP() (*net.UDPConn, error) {
	ip, port, err := a.GetIPPort()
	if err != nil {
		return nil, err
	}
	addr := net.UDPAddr{IP: ip, Port: port}
	return net.ListenUDP(a.Network, &addr)
}

func (a *SimpleAddr) DailUDP() (*net.UDPConn, error) {
	ip, port, err := a.GetIPPort()
	if err != nil {
		return nil, err
	}
	addr := net.UDPAddr{IP: ip, Port: port}
	return net.DialUDP(a.Network, nil, &addr)
}

func (a *SimpleAddr) Listen() (net.Listener, error) {
	switch a.Network {
	case "tcp":
		return net.Listen(a.Network, a.Addr)
	case "unix":
		// 1. remove socket file
		if _, err := os.Stat(a.Addr); err == nil {
			if err := os.Remove(a.Addr); err != nil {
				return nil, err
			}
		}
		// 2. create addr
		addr := &net.UnixAddr{Name: a.Addr, Net: "unix"}
		return net.ListenUnix(a.Network, addr)
	default:
		return nil, fmt.Errorf(`not support netwrok: %s`, a.Network)
	}
}

func (a *SimpleAddr) Dial() (net.Conn, error) {
	switch a.Network {
	case "tcp":
		return net.Dial(a.Network, a.Addr)
	case "unix":
		addr := &net.UnixAddr{Name: a.Addr, Net: "unix"}
		return net.DialUnix(a.Network, nil, addr)
	default:
		return nil, fmt.Errorf(`not support netwrok: %s`, a.Network)
	}
}

func MustParseAddr(addr string) *SimpleAddr {
	parseAddr, err := ParseAddr(addr)
	if err != nil {
		panic(err)
	}
	return parseAddr
}

func ParseAddr(addr string) (*SimpleAddr, error) {
	if addr == "" {
		return nil, fmt.Errorf(`invalid addr: %s`, addr)
	}
	ret := &SimpleAddr{Network: "tcp"}
	splits := strings.SplitN(addr, "://", 2)
	if len(splits) == 1 { // not contains ://
		if strings.Contains(addr, ".socket") || strings.Contains(addr, ".sock") {
			ret.Network = "unix"
		}
		ret.Addr = addr
		return ret, nil
	}
	ret.Network = splits[0]
	ret.Addr = splits[1]
	if !allowedNetwork[ret.Network] {
		return nil, fmt.Errorf(`invalid network: %s`, ret.Network)
	}
	if ret.Addr == "" {
		return nil, fmt.Errorf(`invalid addr: %s`, addr)
	}
	return ret, nil
}
