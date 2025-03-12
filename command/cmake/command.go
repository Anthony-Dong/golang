package cmake

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/command/cmake/static"
	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "cmake",
		Short: "cmake tools",
		Example: `
1. devtool cmake init --name xxx
2. cmake -DCMAKE_BUILD_TYPE=Debug -DCMAKE_EXPORT_COMPILE_COMMANDS=TRUE -S . -B output
3. cmake --build output --config Debug --target main -j 8
4. ./output/main
`,
	}
	if err := command.AddCommand(cmd, NewInitCommand); err != nil {
		return nil, err
	}
	return cmd, nil
}

func NewInitCommand() (*cobra.Command, error) {
	name := ""
	dir := ""
	cmd := &cobra.Command{
		Use:   "init",
		Short: "init cmake project",
		RunE: func(cmd *cobra.Command, args []string) error {
			files, err := static.ReadAllFiles()
			if err != nil {
				return err
			}
			for file, content := range files {
				if file == static.CMakeListsFile {
					content = strings.ReplaceAll(content, "PROJECT_NAME", name)
				}
				if err := utils.WriteFileForce(filepath.Join(dir, file), []byte(content)); err != nil {
					return err
				}
			}
			logs.Info("init cmake project success.")
			logs.Info("step1(config): cmake -DCMAKE_BUILD_TYPE=Debug -DCMAKE_EXPORT_COMPILE_COMMANDS=TRUE -S . -B output")
			logs.Info("step2(build): cmake --build output --config Debug --target main -j 8")
			logs.Info("step3(run): ./output/main")
			return nil
		},
	}
	cmd.Flags().StringVarP(&name, "name", "N", "", "the project name")
	if err := cmd.MarkFlagRequired("name"); err != nil {
		return nil, err
	}
	cmd.Flags().StringVarP(&dir, "dir", "D", utils.GetPwd(), "the project dir")
	return cmd, nil
}
