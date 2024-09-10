package codec

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/codec/thrift_codec"
)

// echo "AAAAEYIhAQRUZXN0HBwWAhUCAAAA" | bin/gtool codec base64 --decode | bin/gtool codec thrift | jq
func newThriftCodecCmd() (*cobra.Command, error) {
	debug := false
	debugLog := func(format string, v ...interface{}) {
		if !debug {
			return
		}
		fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", v...)
	}
	cmd := &cobra.Command{
		Use:     "thrift",
		Short:   "decode thrift protocol binary message",
		Example: `	echo "AAAAEYIhAQRUZXN0HBwWAhUCAAAA" | devtool codec base64 --decode | devtool codec thrift | jq`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !utils.CheckStdInFromPiped() {
				return cmd.Help()
			}
			ctx := cmd.Context()
			input, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf(`read Stdin find err: %v`, err)
			}

			if decodeString, err := base64.StdEncoding.DecodeString(strings.TrimSpace(string(input))); err == nil {
				debugLog("the input data encode type is base64")
				input = decodeString
			}

			handlerStruct := func(payload []byte, proto thrift_codec.Protocol) error {
				data, err := thrift_codec.DecodeStruct(ctx, thrift_codec.NewTProtocol(bytes.NewBuffer(payload), proto))
				if err != nil {
					return fmt.Errorf(`decode struct find err(proto=%s): %v`, proto, err)
				}
				_, _ = os.Stdout.WriteString(utils.ToJson(data))
				return nil
			}

			handlerMessage := func(payload []byte) error {
				buffer := bufio.NewReader(bytes.NewBuffer(payload))
				protocol, err := thrift_codec.GetProtocol(ctx, buffer)
				if err != nil {
					return fmt.Errorf(`decode message find err: %v`, err)
				}
				data, err := thrift_codec.DecodeMessage(ctx, thrift_codec.NewTProtocol(buffer, protocol))
				if err != nil {
					return fmt.Errorf(`decode message find err(proto=%s): %v`, protocol, err)
				}
				data.Protocol = protocol
				_, _ = os.Stdout.WriteString(utils.ToJson(data))
				return nil
			}

			if err = handlerMessage(input); err == nil {
				return nil
			}
			debugLog("handlerMessage find err: %v", err.Error())

			for _, proto := range []thrift_codec.Protocol{
				thrift_codec.UnframedBinary,
				thrift_codec.UnframedCompact,
				thrift_codec.FramedBinary,
				thrift_codec.FramedCompact,
			} {
				if err = handlerStruct(input, proto); err == nil {
					debugLog("proto is: %s", proto)
					return nil
				}
				debugLog("handlerStruct(%s) find err: %v", proto, err.Error())
			}

			return fmt.Errorf(`invalid thrift payload`)
		},
	}
	cmd.Flags().BoolVar(&debug, "debug", false, "enable debug")
	return cmd, nil
}
