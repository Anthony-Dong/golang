package cpp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewCommand() (*cobra.Command, error) {
	tools := Tools{
		Dir:         utils.GetPwd(),
		CXX:         CXX(),
		CC:          CC(),
		CompileType: CompileTypeDebug,
	}
	isRun := false
	isRelease := false
	linkType := "binary"
	output := ""
	thread := 1
	cmd := &cobra.Command{
		Use:   "cpp [--src .cpp] [--hdr .h] [-o binary] [--type binary] [--thread number] [-r] [flags] -- [build flags ... ]",
		Short: "The cpp language tools",
		Long:  "Supports fast compile and running of a cpp file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if isRelease {
				tools.CompileType = CompileTypeRelease
			}
			if !utils.ExistDir(filepath.Join(tools.Dir, ToolsOutputDir)) {
				if err := os.Mkdir(filepath.Join(tools.Dir, ToolsOutputDir), utils.DefaultDirMode); err != nil {
					return err
				}
			}
			if output == "" {
				if linkType == LinkTypeBinary {
					mainFile := filepath.Base(tools.SRCS[len(tools.SRCS)-1])
					output = strings.TrimSuffix(mainFile, filepath.Ext(mainFile))
				} else {
					return fmt.Errorf(`required output flag`)
				}
			}
			if linkType != LinkTypeLibrary && linkType != LinkTypeBinary {
				return fmt.Errorf(`not support link type: %s`, linkType)
			}
			tools.BuildArgs = args
			logs.Debug("output: %s, link type: %s, thread number: %d, tools config: %s", output, linkType, thread, utils.ToJson(tools, true))
			if err := tools.Init(); err != nil {
				return err
			}
			if err := tools.Compile(thread); err != nil {
				return err
			}
			if err := tools.Link(linkType, output); err != nil {
				return err
			}
			if isRun {
				return tools.Run(output)
			}
			return nil
		},
	}
	cmd.Flags().StringSliceVar(&tools.SRCS, "src", []string{}, "The source files")
	cmd.Flags().StringSliceVar(&tools.HDRS, "hdr", []string{}, "The source header files")
	cmd.Flags().StringVar(&linkType, "type", "binary", "The link object type [binary|library]")
	cmd.Flags().BoolVarP(&isRun, "run", "r", false, "Exec output binary file")
	cmd.Flags().BoolVar(&isRelease, "release", false, "Set the compile type is release")
	cmd.Flags().StringVarP(&output, "output", "o", "", "The output file")
	cmd.Flags().IntVarP(&thread, "thread", "j", 1, "The number of compiled threads")

	cmd.Flags().StringSliceVarP(&tools.BuildIncludes, "include", "I", []string{}, "Add directory to include search path")
	cmd.Flags().StringSliceVarP(&tools.LinkIncludes, "link_include", "L", []string{}, "Add directory to library search path")
	cmd.Flags().StringSliceVarP(&tools.LinkLibraries, "link", "l", []string{}, "Add link library")
	if err := cmd.MarkFlagRequired("src"); err != nil {
		return nil, err
	}
	return cmd, nil
}
