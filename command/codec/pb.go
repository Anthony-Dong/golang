package codec

import (
	"fmt"
	"io"
	"os"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/codec/pb_codec"
	"github.com/anthony-dong/golang/pkg/codec/pb_codec/codec"
)

func newPBCodecCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     "pb",
		Short:   "Decode Protocol Buffer (Protobuf) binary messages",
		Example: `echo "CgVoZWxsbxCIBEIDCIgE" | devtool codec pb`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !utils.CheckStdInFromPiped() {
				return cmd.Help()
			}
			in, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf(`read std.in find err: %v`, err)
			}
			message, err := pb_codec.DecodeMessage(cmd.Context(), codec.NewBuffer(in))
			if err != nil {
				return fmt.Errorf(`decode pb message find err: %v`, err)
			}
			_, _ = os.Stdout.WriteString(utils.ToJson(message))
			return nil
		},
	}
	return cmd, nil
}
