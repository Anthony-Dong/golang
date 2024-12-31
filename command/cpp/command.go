package cpp

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/utils"
)

func NewCommand() (*cobra.Command, error) {
	tools := Tools{Pwd: utils.GetPwd()}
	configFile := ""
	isRun := false
	linkType := "binary"
	thread := 1
	cmd := &cobra.Command{
		Use:   "cpp [--src .cpp] [--hdr .h] [-o binary] [--type binary] [--thread number] [-r] [flags]",
		Short: "The cpp language tools",
		Long:  "Supports fast compile and running of a cpp file",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if configFile != "" {
				if err := ReadToolConfigFromFile(configFile, &tools); err != nil {
					return err
				}
			}
			if !utils.Contains([]string{LinkTypeBinary, LinkTypeLibrary}, linkType) {
				return fmt.Errorf(`not support link type: %s`, linkType)
			}
			if linkType == LinkTypeLibrary && filepath.Ext(tools.Output) != ".a" {
				return fmt.Errorf(`if the link type is libary then the output file type must be .a`)
			}
			if err := tools.Build(ctx, thread); err != nil {
				return err
			}
			if err := tools.Link(ctx, linkType); err != nil {
				return err
			}
			if isRun && linkType == LinkTypeBinary {
				return tools.Run(ctx)
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&tools.SRCS, "src", []string{}, "The source files")
	cmd.Flags().StringSliceVar(&tools.HDRS, "hdr", []string{}, "The source header files")
	cmd.Flags().StringVar(&linkType, "link_type", LinkTypeBinary, "The link object type [binary|library]")
	cmd.Flags().StringVar(&tools.BuildType, "build_type", BuildTypeDebug, "set the build type is release")
	cmd.Flags().StringVarP(&tools.Output, "output", "o", "", "The output file")
	cmd.Flags().IntVarP(&thread, "thread", "j", 1, "The number of compiled threads")
	cmd.Flags().StringSliceVar(&tools.BuildArgs, "cxxopt", []string{}, "The build opt")
	cmd.Flags().StringSliceVar(&tools.LinkArgs, "linkopt", []string{}, "The Link opt")
	cmd.Flags().BoolVarP(&isRun, "run", "r", false, "Exec output binary file")
	cmd.Flags().StringVar(&configFile, "config", "", "set build config")
	if err := cmd.MarkFlagRequired("src"); err != nil {
		return nil, err
	}
	if err := cmd.MarkFlagRequired("output"); err != nil {
		return nil, err
	}
	return cmd, nil
}
