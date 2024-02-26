package middleware

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/pkg/logs"
)

func NewInitLoggerMv(config *command.AppConfig) command.Middleware {
	return func(ctx *cobra.Command, args []string) error {
		if strings.HasPrefix(ctx.Use, "curl") {
			logs.SetPrinterStdError()
		}
		if config.LogLevel != "" {
			logs.SetLevelString(config.LogLevel)
		}
		if config.Verbose {
			logs.SetLevel(logs.LevelDebug)
		}
		return nil
	}
}
