package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/anthony-dong/golang/command/golang"

	"github.com/anthony-dong/golang/command/cpp"
	"github.com/anthony-dong/golang/command/git"

	"github.com/anthony-dong/golang/command/turl"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/command/jsontool"
	"github.com/anthony-dong/golang/command/middleware"
	"github.com/anthony-dong/golang/command/run"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command/codec"
	"github.com/anthony-dong/golang/command/gen"
	"github.com/anthony-dong/golang/command/hexo"
	"github.com/anthony-dong/golang/command/tcpdump"
	"github.com/anthony-dong/golang/command/upload"
	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

func main() {
	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	defer command.Close()

	ctx, cancel := context.WithCancel(context.Background())
	command.Defer(cancel)

	go func() {
		defer close(done)
		cmd, err := NewCmd()
		if err != nil {
			command.Defer(func() { command.ExitError(err) })
			return
		}
		if err := cmd.ExecuteContext(ctx); err != nil {
			command.Defer(func() { command.ExitError(err) })
			return
		}
	}()
	<-done
}

func NewCmd() (*cobra.Command, error) {
	config := &command.AppConfig{}
	var cmd = &cobra.Command{
		Use:                   command.AppName,
		Version:               command.AppVersion,
		CompletionOptions:     cobra.CompletionOptions{DisableDefaultCmd: true},
		SilenceUsage:          true, // 禁止失败打印 --help
		SilenceErrors:         true, // 禁止框架打印异常
		DisableFlagsInUseLine: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if err := middleware.NewInitLoggerMv(config)(cmd, args); err != nil {
				return err
			}
			if err := middleware.NewInitConfigMv(config)(cmd, args); err != nil {
				return err
			}
			logs.Debug("start cmd: %s, cmd.args: %s, os.args: %s", command.AppName, utils.ToJson(args), utils.ToJson(os.Args))
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
	if err := command.AddConfigCommand(cmd, config, hexo.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddConfigCommand(cmd, config, upload.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, gen.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, tcpdump.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddConfigCommand(cmd, config, run.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, golang.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, turl.NewTurlCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, git.NewCommand); err != nil {
		return nil, err
	}
	if err := command.AddCommand(cmd, cpp.NewCommand); err != nil {
		return nil, err
	}
	return cmd, nil
}
