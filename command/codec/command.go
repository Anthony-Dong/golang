package codec

import (
	"fmt"
	"io"
	"os"

	"github.com/anthony-dong/golang/command"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/codec"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "codec",
		Short: "The Encode and Decode data tool",
	}
	if err := command.AddCommand(cmd, newThriftCodecCmd); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, newPBCodecCmd); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewProtocCodec); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewJson2YamlCmd); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewYaml2JsonCmd); err != nil {
		return nil, err
	}
	cmd.AddCommand(newCodecCmd("gzip", codec.NewGzipCodec()))
	cmd.AddCommand(newCodecCmd("base64", codec.NewCodec(codec.NewBase64Codec())))
	cmd.AddCommand(newCodecCmd("br", codec.NewBrCodec()))
	cmd.AddCommand(newCodecCmd("deflate", codec.NewDeflateCodec()))
	cmd.AddCommand(newCodecCmd("snappy", codec.NewCodec(codec.NewSnappyCodec())))
	cmd.AddCommand(newCodecCmd("md5", codec.NewCodec(codec.NewMd5Codec())))
	cmd.AddCommand(newCodecCmd("hex", codec.NewCodec(codec.NewHexCodec())))
	cmd.AddCommand(newCodecCmd("hexdump", codec.NewCodec(codec.NewHexDumpCodec())))
	cmd.AddCommand(newCodecCmd("pb-desc", NewDecoderFunc(codec.NewProtoDesc().Decode)))
	cmd.AddCommand(newCodecCmd("url", NewDecoderFunc(ParseUrl)))
	cmd.AddCommand(newCodecCmd("cmd", NewDecoderFunc(ParseCmd)))
	return cmd, nil
}

func newCodecCmd(name string, codec codec.Codec) *cobra.Command {
	var (
		reader   io.Reader = os.Stdin
		writer   io.Writer = os.Stdout
		isDecode bool
		isLf     bool
	)
	cmd := &cobra.Command{
		Use:   name,
		Short: fmt.Sprintf("%s codec", name),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if !utils.CheckStdInFromPiped() {
				return cmd.Help()
			}
			defer func() {
				if err == nil && isLf {
					_, _ = writer.Write([]byte{'\n'})
				}
			}()
			if isDecode {
				return codec.Decode(reader, writer)
			}
			return codec.Encode(reader, writer)
		},
	}
	cmd.Flags().BoolVar(&isLf, "lf", false, "append lf")
	cmd.Flags().BoolVar(&isDecode, "decode", false, "decode content data")
	return cmd
}

func NewDecoderFunc(fun BytesDecoderFunc) codec.Codec {
	return codec.NewCodec(fun)
}

func NewEncodeFunc(fun BytesEncodeFunc) codec.Codec {
	return codec.NewCodec(fun)
}

type BytesDecoderFunc func(src []byte) (dst []byte, err error)
type BytesEncodeFunc func(src []byte) (dst []byte)

func (b BytesDecoderFunc) Decode(src []byte) (dst []byte, err error) {
	return b(src)
}
func (b BytesDecoderFunc) Encode(src []byte) (dst []byte) {
	return src
}

func (b BytesEncodeFunc) Encode(src []byte) (dst []byte) {
	return b(src)
}

func (b BytesEncodeFunc) Decode(src []byte) (dst []byte, err error) {
	return nil, fmt.Errorf(`not support decode`)
}
