package middleware

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/pkg/logs"
)

func NewInitLoggerMv(verbose bool, logLevel string) command.Middleware {
	return func(ctx *cobra.Command, args []string) error {
		if strings.HasPrefix(ctx.Use, "curl") {
			logs.SetPrinterStdError()
		}
		if logLevel != "" {
			logs.SetLevelString(logLevel)
		}
		if verbose {
			logs.SetLevel(logs.LevelDebug)
		}
		return nil
	}
}
