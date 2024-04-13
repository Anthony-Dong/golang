package cli

import (
	"context"

	"github.com/anthony-dong/golang/command/proxy"
	"github.com/anthony-dong/golang/command/tcp"

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
	"github.com/anthony-dong/golang/command/run"
	"github.com/anthony-dong/golang/command/tcpdump"
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
			if err := middleware.NewInitLoggerMv(config)(cmd, args); err != nil {
				return err
			}
			if err := middleware.NewInitConfigMv(config)(cmd, args); err != nil {
				return err
			}
			cmd.SetContext(context.WithValue(ctx, command.AppConfigCtxKey, config))
			return nil
		},
	}
	cmd.SetHelpTemplate(command.HelpTmpl)
	cmd.SetUsageTemplate(command.UsageTmpl)
	cmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "Turn on verbose mode")
	cmd.PersistentFlags().StringVar(&config.ConfigFile, "config-file", config.ConfigFile, "Set the config file")
	cmd.PersistentFlags().StringVar(&config.LogLevel, "log-level", "", "Set the log level in [debug|info|notice|warn|error] (default \"info\")")
	if err := command.AddCommand(cmd, codec.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, jsontool.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, hexo.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, upload.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, gen.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, tcpdump.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, run.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, golang.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, curl.NewCurlCommand); err != nil {
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
	return cmd, nil
}
