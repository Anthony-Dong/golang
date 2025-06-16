package codegen

import (
	"os"
	"path/filepath"

	"github.com/anthony-dong/golang/command"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/logs"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{Use: "codegen", Short: "Generate code from IDL files (Thrift, Protobuf)"}
	if err := command.AddCommand(cmd, NewProtocCommand); err != nil {
		return nil, err
	}
	return cmd, nil
}

func NewProtocCommand() (*cobra.Command, error) {
	var (
		dir     string
		gopkg   string
		include []string
		output  string
	)
	cmd := cobra.Command{
		Use:   "protoc [-D dir] [-I include] [--go_pkg package]",
		Short: "Compile Protocol Buffer (Protobuf) IDL files",
		Long: `Plugin make it easy to compile PB:
Help Doc:
	golang: https://go.dev
	proto3: https://developers.google.com/protocol-buffers/docs/proto3
	protoc-gen-go: https://developers.google.com/protocol-buffers/docs/reference/go-generated
	grpc: https://grpc.io/docs/what-is-grpc/introduction
	protoc-gen-go-grpc: https://grpc.io/docs/languages/go/quickstart
Install:
	golang: https://go.dev/dl
	protoc: https://github.com/protocolbuffers/protobuf/releases
	protoc-gen-go: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	protoc-gen-go-grpc: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			gen := NewProtocGen(dir, gopkg, func(gen *ProtocGen) {
				if len(include) > 0 {
					gen.Include = include
				}
				if output != "" {
					gen.OutPutPath = output
				}
			})
			result, err := gen.Gen(ctx)
			if err != nil {
				if result == nil {
					return err
				}
				logs.CtxError(ctx, "[Protoc] exec error\n===STD OUT===\n%s\n===ERR OUT====\n%s", result.StdOut.String(), result.StdError.String())
				return err
			}
			logs.CtxInfo(ctx, "[Protoc] exec success\n%s", result.Command)
			return nil
		},
	}
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	cmd.Flags().StringVarP(&dir, "dir", "D", pwd, "The project dir")
	cmd.Flags().StringVarP(&gopkg, "go_pkg", "", "", "Define go import path, eg: anthony-dong/proto-example/pb-gen")
	cmd.Flags().StringVarP(&output, "output", "O", filepath.Join(utils.GetGoPath(), "src"), "The output dir")
	cmd.Flags().StringArrayVarP(&include, "include", "I", []string{}, "Add an IDL search path for includes")
	if err := cmd.MarkFlagRequired("go_pkg"); err != nil {
		return nil, err
	}
	return &cmd, nil
}
