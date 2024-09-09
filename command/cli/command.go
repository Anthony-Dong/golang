package cli

import (
	"context"

	"github.com/anthony-dong/golang/command/_init"
	"github.com/anthony-dong/golang/command/install"
	"github.com/anthony-dong/golang/command/mock"

	"github.com/anthony-dong/golang/command/tools"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/command/codec"
	"github.com/anthony-dong/golang/command/cpp"
	"github.com/anthony-dong/golang/command/curl"
	"github.com/anthony-dong/golang/command/gen"
	"github.com/anthony-dong/golang/command/git"
	"github.com/anthony-dong/golang/command/golang"
	"github.com/anthony-dong/golang/command/hexo"
	"github.com/anthony-dong/golang/command/jsontool"
	"github.com/anthony-dong/golang/command/middleware"
	"github.com/anthony-dong/golang/command/proxy"
	"github.com/anthony-dong/golang/command/run"
	"github.com/anthony-dong/golang/command/tcp"
	"github.com/anthony-dong/golang/command/upload"
)

func NewCommand(config *command.AppConfig) (*cobra.Command, error) {
	if config == nil {
		config = &command.AppConfig{}
	}
	if config.AppName == "" {
		config.AppName = command.AppName
	}
	if config.AppVersion == "" {
		config.AppVersion = command.AppVersion
	}
	var (
		verbose    bool
		logLevel   string
		configFile string
	)
	var cmd = &cobra.Command{
		Use:                   config.AppName,
		Version:               config.AppVersion,
		CompletionOptions:     cobra.CompletionOptions{DisableDefaultCmd: true},
		SilenceUsage:          true, // 禁止失败打印 --help
		SilenceErrors:         true, // 禁止框架打印异常
		DisableFlagsInUseLine: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if ctx == nil {
				ctx = context.Background()
			}
			if err := middleware.NewInitLoggerMv(verbose, logLevel)(cmd, args); err != nil {
				return err
			}
			if err := middleware.NewInitConfigMv(configFile, config)(cmd, args); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.SetHelpTemplate(command.HelpTmpl)
	cmd.SetUsageTemplate(command.UsageTmpl)
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Turn on verbose mode")
	cmd.PersistentFlags().StringVar(&configFile, "config-file", "", "Set the config file")
	cmd.PersistentFlags().StringVar(&logLevel, "log-level", "", `Set the log level in [debug|info|notice|warn|error] (default "info")`)
	if err := command.AddCommand(cmd, codec.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, jsontool.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommandWithConfig(cmd, func() *command.HexoConfig { return config.HexoConfig }, hexo.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommandWithConfig(cmd, func() *command.UploadConfig {
		return config.UploadConfig
	}, upload.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, gen.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommandWithConfig(cmd, func() *command.RunTaskConfig {
		return config.RunTaskConfig
	}, run.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, golang.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommandWithConfig(cmd, func() *command.CurlConfig {
		return config.CurlConfig
	}, curl.NewCurlCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, git.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, cpp.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, proxy.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, tcp.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, tools.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, _init.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, mock.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, install.NewCommand); err != nil {
		return nil, err
	}
	return cmd, nil
}
