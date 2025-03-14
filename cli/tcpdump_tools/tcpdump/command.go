package tcpdump

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"github.com/anthony-dong/golang/pkg/logs"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/codec"
	"github.com/anthony-dong/golang/pkg/tcpdump"
	"github.com/anthony-dong/golang/pkg/utils"
)

func joinCommand(name string) string {
	cmd := filepath.Base(os.Args[0])
	if cmd == name {
		return name
	}
	return cmd + " " + name
}

func NewCommand(name string, cmdOption func(command *cobra.Command), ops ...DecodeOption) (*cobra.Command, error) {
	var (
		filename       string
		verbose        bool
		filterMsgTypes []string
	)
	cmd := &cobra.Command{
		Use:   fmt.Sprintf(`%s [-r file] [-v]`, name),
		Short: `Decode tcpdump packets`,
		Example: fmt.Sprintf(`  
1. 在线抓包: tcpdump 'port 8888' -X -l -n | %s
2. 离线分析:
 - 抓取流量: tcpdump 'port 8888' -w xxx.pcap
 - 分析流量: %[1]s -r xxx.pcap`, joinCommand(name)),
		RunE: func(cmd *cobra.Command, args []string) error {
			option := NewDecodeOptions(ops...)
			if option.MsgWriter == nil {
				if len(filterMsgTypes) == 0 {
					filterMsgTypes = append(filterMsgTypes, tcpdump.MessageType_Thrift, tcpdump.MessageType_HTTP)
					switch GetPacketSourceType(filename) {
					case PacketSource_Consul:
						filterMsgTypes = append(filterMsgTypes, tcpdump.MessageType_Tcpdump)
					case PacketSource_File:
						filterMsgTypes = append(filterMsgTypes, tcpdump.MessageType_TcpPacket)
					}
					if verbose {
						filterMsgTypes = append(filterMsgTypes, tcpdump.MessageType_Log)
					}
				}
				logs.CtxInfo(cmd.Context(), "filter msg type: %s", utils.ToJson(filterMsgTypes))
				option.MsgWriter = tcpdump.NewConsoleLogMessageWriter(tcpdump.ConvertToMessageType(filterMsgTypes))
			}
			if len(option.Decoders) == 0 {
				option.Decoders = append(option.Decoders, tcpdump.NewThriftDecoder(), tcpdump.NewHttpDecoder())
			}
			source, err := NewPacketSource(filename, option)
			if err != nil {
				return err
			}
			return DecodePacketSource(cmd.Context(), source, option)
		},
	}
	if cmdOption != nil {
		cmdOption(cmd)
	}
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "开启DEBUG模式, 会打印一些DEBUG信息")
	cmd.Flags().StringVarP(&filename, "file", "r", "", "读取tcpdump抓取的pcap文件")
	cmd.Flags().StringArrayVarP(&filterMsgTypes, "filter", "f", []string{}, "过滤消息类型(thrift/http/log/tcpdump)")
	return cmd, nil
}

func NewPacketSource(filename string, options DecodeOptions) (PacketSource, error) {
	if filename != "" {
		return NewFileSource(filename, options)
	}
	if utils.CheckStdInFromPiped() {
		source := NewConsulSource(os.Stdin, options)
		return source, nil
	}
	return nil, fmt.Errorf("no packet source")
}

func GetPacketSourceType(filename string) PacketSourceType {
	if filename != "" {
		return PacketSource_File
	}
	if utils.CheckStdInFromPiped() {
		return PacketSource_Consul
	}
	return PacketSource_Unknown
}

func DecodePacketSource(ctx context.Context, source PacketSource, option DecodeOptions) error {
	decoder := tcpdump.NewPacketDecoder(option.MsgWriter)
	for _, v := range option.Decoders {
		decoder.AddDecoder(v)
		logs.CtxInfo(ctx, "use decoder %s", v.Name())
	}
	for data := range source.Packets() {
		decoder.Decode(ctx, NewTcpPacket(data, option.MsgWriter))
		if wait, isOk := data.(WaitPacket); isOk {
			wait.Notify()
		}
	}
	return nil
}

var packetCounter = 1

func NewTcpPacket(packet gopacket.Packet, writer tcpdump.MessageWriter) *tcpdump.TcpPacket {
	var (
		src, dest         net.IP
		L3IsOk, L4IsOK    bool
		srcPort, destPort int
		tcpFlags          []string
		data              = &tcpdump.TcpPacket{}
	)
	switch L3 := packet.NetworkLayer().(type) {
	case *layers.IPv4:
		L3IsOk = true
		src = L3.SrcIP
		dest = L3.DstIP
	case *layers.IPv6:
		L3IsOk = true
		src = L3.SrcIP
		dest = L3.DstIP
	}
	switch L4 := packet.TransportLayer().(type) {
	case *layers.TCP:
		L4IsOK = true
		srcPort = int(L4.SrcPort)
		destPort = int(L4.DstPort)
		tcpFlags = GetTcpFlags(L4)
	}
	if L3IsOk && L4IsOK {
		data.Src = tcpdump.IpPort(src.String(), srcPort)
		data.Dst = tcpdump.IpPort(dest.String(), destPort)
		data.TCPFlag = tcpFlags
		tcp := packet.TransportLayer().(*layers.TCP)
		result := HandlerTcp(data.Src, data.Dst, tcp)
		if !result.Is(OutOfOrderStatus) {
			data.Data = tcp.Payload
		}
		data.ACK = int(tcp.Ack)
		header := strings.Builder{}
		header.WriteString(fmt.Sprintf("[%d] ", packetCounter))
		header.WriteString(fmt.Sprintf("[%s] ", packet.Metadata().Timestamp.Format(utils.FormatTimeV1)))
		header.WriteString(fmt.Sprintf("[%s-%s] ", packet.NetworkLayer().LayerType(), packet.TransportLayer().LayerType()))
		header.WriteString(fmt.Sprintf("[%s -> %s] ", data.Src, data.Dst))
		header.WriteString(fmt.Sprintf("[%s] ", strings.Join(tcpFlags, ",")))
		header.WriteString(fmt.Sprintf("%s ", GetRelativeInfo(data.Src, data.Dst, tcp)))
		if len(tcp.Payload) != 0 {
			header.WriteString(fmt.Sprintf("[%d Byte] ", len(tcp.Payload)))
		}
		header.WriteString(fmt.Sprintf("%v", result))
		data.Header = header.String()
		packetCounter = packetCounter + 1
		return data
	}
	if layer := packet.TransportLayer(); layer != nil {
		writer.Write(&TransportMessage{Layer: layer})
	}
	return data
}

type TransportMessage struct {
	Layer gopacket.TransportLayer
}

func (m *TransportMessage) String() string {
	return string(codec.NewHexDumpCodec().Encode(m.Layer.LayerPayload()))
}

func (*TransportMessage) Type() tcpdump.MessageType {
	return tcpdump.MessageType_Layer
}

func GetTcpFlags(L4 *layers.TCP) []string {
	var flags []string
	if L4.FIN {
		flags = append(flags, "FIN")
	}
	if L4.SYN {
		flags = append(flags, "SYN")
	}
	if L4.ACK {
		flags = append(flags, "ACK")
	}
	if L4.PSH {
		flags = append(flags, "PSH")
	}
	if L4.RST {
		flags = append(flags, "RST")
	}
	if L4.URG {
		flags = append(flags, "URG")
	}
	if L4.ECE {
		flags = append(flags, "ECE")
	}
	if L4.CWR {
		flags = append(flags, "CWR")
	}
	if L4.NS {
		flags = append(flags, "NS")
	}
	return flags
}
