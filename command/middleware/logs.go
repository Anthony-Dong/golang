package middleware

import (
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/pkg/logs"
)

type Middleware func(cmd *cobra.Command, args []string) error

func NewInitLoggerMv(config *command.AppConfig) Middleware {
	return func(ctx *cobra.Command, args []string) error {
		if config.LogLevel != "" {
			logs.SetLevelString(config.LogLevel)
		}
		if config.Verbose {
			logs.SetLevel(logs.LevelDebug)
		}
		return nil
	}
}
