package codec

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"

	"github.com/anthony-dong/golang/pkg/utils"
)

func NewProtocCodec() (*cobra.Command, error) {
	isMockProto := false
	cmd := &cobra.Command{
		Use:   "protoc",
		Short: `protoc desc codec tools`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !utils.CheckStdInFromPiped() {
				return cmd.Help()
			}
			stdIn, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf(`read stdin find err: %v`, err)
			}
			data := descriptorpb.FileDescriptorSet{}
			if err := proto.Unmarshal(stdIn, &data); err != nil {
				return err
			}
			json, err := protoMessageToJson(&data)
			if err != nil {
				return err
			}
			if _, err := os.Stdout.Write(json); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&isMockProto, "mock", false, "mock protoc plugin")
	return cmd, nil
}

func protoMessageToJson(m proto.Message) ([]byte, error) {
	return protojson.MarshalOptions{Multiline: true, AllowPartial: true}.Marshal(m)
}
