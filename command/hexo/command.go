package hexo

import (
	"path/filepath"

	"github.com/anthony-dong/golang/command"

	"github.com/anthony-dong/golang/pkg/logs"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/utils"
)

func NewCommand(config *command.AppConfig) (*cobra.Command, error) {
	cmd := &cobra.Command{Use: "hexo", Short: "The Hexo tool"}
	if err := command.AddConfigCommand(cmd, config, NewBuildCmd); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, NewReadmeCmd); err != nil {
		return nil, err
	}
	return cmd, nil
}

type hexoConfig struct {
	Dir       string   `json:"dir"`
	Keyword   []string `json:"keyword"`
	Ignore    []string `json:"ignore"`
	TargetDir string   `json:"target_dir"`
}

func NewBuildCmd(config *command.AppConfig) (*cobra.Command, error) {
	var (
		cfg = &hexoConfig{}
	)
	cmd := &cobra.Command{Use: "build", Short: "Build the markdown project to hexo"}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if config.HexoConfig == nil {
			config.HexoConfig = &command.HexoConfig{}
		}
		for _, elem := range config.HexoConfig.KeyWord {
			cfg.Keyword = append(cfg.Keyword, elem)
		}
		for _, elem := range config.HexoConfig.Ignore {
			cfg.Ignore = append(cfg.Ignore, elem)
		}
		if dir, err := filepath.Abs(cfg.Dir); err != nil {
			return err
		} else {
			cfg.Dir = dir
		}
		if dir, err := filepath.Abs(cfg.TargetDir); err != nil {
			return err
		} else {
			cfg.TargetDir = dir
		}
		logs.Info("[Hexo] init config success:\n%s", utils.ToJson(cfg, true))
		return buildHexo(cmd.Context(), cfg.Dir, cfg.TargetDir, cfg.Keyword, cfg.Ignore)
	}
	cmd.Flags().StringVarP(&cfg.Dir, "dir", "d", "", "The markdown project dir (required)")
	cmd.Flags().StringVarP(&cfg.TargetDir, "target_dir", "t", "", "The hexo post dir (required)")
	cmd.Flags().StringArrayVarP(&cfg.Keyword, "keyword", "k", nil, "The keyword need clear, eg: baidu => xxxxx, read from command and load config")
	if err := cmd.MarkFlagRequired("dir"); err != nil {
		return nil, err
	}
	if err := cmd.MarkFlagRequired("target_dir"); err != nil {
		return nil, err
	}
	return cmd, nil
}
