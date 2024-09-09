package install

import (
	"context"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command/install/static"
	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewCommand() (*cobra.Command, error) {
	filename := ""
	list := false
	cmd := &cobra.Command{Use: "install", Short: "install script", RunE: func(cmd *cobra.Command, args []string) error {
		if list {
			for _, file := range static.GetFiles() {
				logs.CtxInfo(cmd.Context(), "file: %s", file)
			}
			return nil
		}
		if filename == "" {
			return cmd.Usage()
		}
		return install(cmd.Context(), filename)
	}}
	cmd.Flags().StringVar(&filename, "file", "", "install file")
	cmd.Flags().BoolVar(&list, "list", false, "list files")
	return cmd, nil
}

func install(ctx context.Context, filename string) error {
	dst := filepath.Join(utils.GetUserHomeDir(), "go", "bin")
	src, err := static.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dst, filename), src, 0777); err != nil {
		return err
	}
	logs.CtxInfo(ctx, "install script [%s] success", filename)
	files, err := static.GetExtraFiles(filename)
	if err != nil {
		return err
	}
	for extraFile, content := range files {
		if err := os.WriteFile(filepath.Join(dst, extraFile), content, 0777); err != nil {
			return err
		}
		logs.CtxInfo(ctx, "install extra file [%s] success", extraFile)
	}
	usage := static.GetUsage(filename)
	if usage != "" {
		logs.CtxInfo(ctx, "usage:\n", usage)
	}
	return nil
}
