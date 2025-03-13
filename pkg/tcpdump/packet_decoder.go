package tcpdump

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"strconv"

	"github.com/anthony-dong/golang/pkg/logs"
)

type Reader interface {
	io.Reader
	Peek(int) ([]byte, error)
}

type Decoder interface {
	Decode(ctx context.Context, r Reader, _ *TcpPacket) (Message, error)
	Name() string
}

type TcpPacket struct {
	Src     string
	Dst     string
	Data    []byte
	ACK     int
	TCPFlag []string

	Header string // Header
}

func (p *TcpPacket) Type() MessageType {
	return MessageType_TcpPacket
}

func (p *TcpPacket) String() string {
	return p.Header
}

func (p *TcpPacket) IsFin() bool {
	for index, elem := range p.TCPFlag {
		if elem == "FIN" && index == 0 {
			return true
		}
	}
	return false
}

func (p *TcpPacket) IsACK() bool {
	for _, elem := range p.TCPFlag {
		if elem == "ACK" {
			return true
		}
	}
	return false
}

type PacketDecoder struct {
	packets map[string] /*src:dest*/ map[int][]byte

	decoder []Decoder
	writer  MessageWriter
}

func NewPacketDecoder(w MessageWriter) *PacketDecoder {
	return &PacketDecoder{
		packets: make(map[string]map[int][]byte),
		writer:  w,
	}
}

func (c *PacketDecoder) AddDecoder(decoder Decoder) {
	c.decoder = append(c.decoder, decoder)
}

func (c *PacketDecoder) closePacket(packet TcpPacket) {
	if c.packets == nil {
		c.packets = map[string]map[int][]byte{}
	}
	delete(c.packets, packet.Src+":"+packet.Dst)
}

func (c *PacketDecoder) findNext(p *TcpPacket) bool {
	key := p.Src + "|" + p.Dst
	if c.packets[key] == nil {
		return false
	}
	for ack := range c.packets[key] {
		if ack > p.ACK {
			return true
		}
	}
	return false
}

func (c *PacketDecoder) Decode(ctx context.Context, p *TcpPacket) {
	if p == nil {
		return
	}
	c.writeMessage(p)
	if c.packets == nil {
		c.packets = map[string]map[int][]byte{}
	}
	key := p.Src + "|" + p.Dst
	if c.packets[key] == nil {
		c.packets[key] = map[int][]byte{}
	}
	if c.packets[key][p.ACK] == nil {
		c.packets[key][p.ACK] = []byte{}
	}

	if p.IsACK() && len(p.Data) > 0 {
		payload := c.packets[key][p.ACK]
		payload = append(payload, p.Data...)
		c.packets[key][p.ACK] = payload

		c.decode(ctx, p, payload, func() {
			delete(c.packets[key], p.ACK)
		})
	}
}

func (c *PacketDecoder) decode(ctx context.Context, packet *TcpPacket, payload []byte, success func()) {
	decode := func(decoder Decoder) bool {
		reader := bufio.NewReader(bytes.NewBuffer(payload))
		msg, err := decoder.Decode(ctx, reader, packet)
		if err != nil {
			c.writeMessage(NewLogMessage(ctx, logs.LevelWarn, "decoder [%s] decode msg find err: %v", decoder.Name(), err))
			return false
		}
		c.writeMessage(msg)
		return true
	}
	for _, decoder := range c.decoder {
		if decode(decoder) {
			success()
			break
		}
	}
}

func (c *PacketDecoder) writeMessage(msg Message) {
	c.writer.Write(msg)
}

// IpPort 支持 ipv6:port, [ipv6]:port, ip:port
func IpPort(ip string, port int) string {
	if IsIPV6(ip) {
		return "[" + ip + "]:" + strconv.Itoa(port)
	}
	return ip + ":" + strconv.Itoa(port)
}

// IsIPV6 支持ipv6, 不支持 [ipv6]
func IsIPV6(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			return false
		case ':':
			return net.ParseIP(s) != nil
		}
	}
	return false
}
