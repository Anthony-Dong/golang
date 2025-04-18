package tcpdump

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"sync"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/pkg/errors"

	"github.com/anthony-dong/golang/pkg/tcpdump"

	"github.com/anthony-dong/golang/pkg/codec"
)

type PacketSource interface {
	Packets() chan gopacket.Packet
}

type DecodeOptions struct {
	gopacket.DecodeOptions
	MsgWriter tcpdump.MessageWriter
	Decoders  []tcpdump.Decoder
}

type DecodeOption func(options *DecodeOptions)

func WithMsgWriter(writer tcpdump.MessageWriter) DecodeOption {
	return func(options *DecodeOptions) {
		options.MsgWriter = writer
	}
}

func WithDecoder(decoder tcpdump.Decoder) DecodeOption {
	return func(options *DecodeOptions) {
		options.Decoders = append(options.Decoders, decoder)
	}
}

type PacketSourceType string

const (
	PacketSource_File    PacketSourceType = "File"
	PacketSource_Consul  PacketSourceType = "Consul"
	PacketSource_Unknown PacketSourceType = "Unknown"
)

func NewDecodeOptions(ops ...DecodeOption) DecodeOptions {
	options := DecodeOptions{
		DecodeOptions: gopacket.DecodeOptions{
			Lazy:                     false,
			NoCopy:                   true,
			DecodeStreamsAsDatagrams: true,
		},
	}
	for _, op := range ops {
		op(&options)
	}
	return options
}

func NewFileSource(file string, cfg DecodeOptions) (PacketSource, error) {
	if file == "" {
		return nil, errors.New(`required file`)
	}
	filename, err := filepath.Abs(file)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("open %s file err", filename))
	}
	src, err := pcap.OpenOffline(filename)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("open %s file err", filename))
	}
	source := gopacket.NewPacketSource(src, layers.LayerTypeEthernet)
	source.DecodeOptions = cfg.DecodeOptions
	return source, nil
}

type ConsulSource struct {
	r       io.Reader
	lines   chan []string
	packets chan gopacket.Packet
	gopacket.DecodeOptions
	initOnce sync.Once
	writer   tcpdump.MessageWriter
}

func NewConsulSource(r io.Reader, cfg DecodeOptions) *ConsulSource {
	return &ConsulSource{
		lines:         make(chan []string, 1),
		packets:       make(chan gopacket.Packet, 1),
		r:             r,
		DecodeOptions: cfg.DecodeOptions,
		writer:        cfg.MsgWriter,
	}
}

func (t *ConsulSource) Read() {
	scanner := bufio.NewScanner(t.r)
	data := make([]string, 0)
	sendLine := func() {
		if len(data) == 0 {
			return
		}
		c := make([]string, len(data))
		copy(c, data)
		t.lines <- c
		data = data[:0]
	}
	for scanner.Scan() && scanner.Err() == nil {
		line := scanner.Text()
		hex, isEnd := codec.ReadHexdump(line)
		if hex != "" {
			data = append(data, hex)
			if isEnd {
				sendLine()
			}
			continue
		}
		sendLine()
		data = append(data, line)
	}
	sendLine()
	close(t.lines)
}

func (t *ConsulSource) init() {
	go func() {
		t.Read()
	}()
	go func() {
		for elem := range t.lines {
			header := elem[0]
			t.writer.Write(&tcpdump.TcpdumpHeader{
				Header: header,
			})
			payload := &bytes.Buffer{}
			for _, line := range elem[1:] {
				payload.WriteString(line)
			}
			if len(payload.Bytes()) == 0 {
				continue
			}
			decode, err := codec.NewHexCodec().Decode(payload.Bytes())
			if err != nil {
				continue
			}
			var wp CustomPacket
			if packet, _ := t.selectIPV4(decode); packet != nil {
				wp = NewCustomPacket(packet)
				t.packets <- wp
				wp.Wait()
				continue
			}
			if packet, _ := t.selectIPV6(decode); packet != nil {
				wp = NewCustomPacket(packet)
				t.packets <- wp
				wp.Wait()
				continue
			}
			t.writer.Write(&tcpdump.TcpdumpPayload{Payload: decode})
		}
		close(t.packets)
	}()
}

func (t *ConsulSource) Packets() chan gopacket.Packet {
	t.initOnce.Do(t.init)
	return t.packets
}

func (t *ConsulSource) selectIPV4(data []byte) (gopacket.Packet, error) {
	packet := gopacket.NewPacket(data, layers.LayerTypeIPv4, t.DecodeOptions)
	if _, isOK := packet.NetworkLayer().(*layers.IPv4); isOK {
		if err := packet.ErrorLayer(); err != nil {
			return nil, err.Error()
		}
		return packet, nil
	}
	return nil, fmt.Errorf(`can not parse as ipv4`)
}

func (t *ConsulSource) selectIPV6(data []byte) (gopacket.Packet, error) {
	packet := gopacket.NewPacket(data, layers.LayerTypeIPv6, t.DecodeOptions)
	if _, isOK := packet.NetworkLayer().(*layers.IPv6); isOK {
		if err := packet.ErrorLayer(); err != nil {
			return nil, err.Error()
		}
		return packet, nil
	}
	return nil, fmt.Errorf(`can not parse as ipv6`)
}

type customPacket struct {
	gopacket.Packet
	notify chan struct{}
}

type WaitPacket interface {
	Notify()
	Wait()
}

type CustomPacket interface {
	WaitPacket
	gopacket.Packet
}

func NewCustomPacket(data gopacket.Packet) CustomPacket {
	return &customPacket{
		notify: make(chan struct{}, 0),
		Packet: data,
	}
}

func (w *customPacket) Notify() {
	close(w.notify)
}

func (w *customPacket) Wait() {
	<-w.notify
}
