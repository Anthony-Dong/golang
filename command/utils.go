package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/logs"
)

func ExitError(err error) {
	if err != nil {
		logs.StdError("[%s] exit error: %s", AppName, err.Error())
	}
	os.Exit(1)
}

func AddCommand(cmd *cobra.Command, foo func() (*cobra.Command, error)) error {
	subCmd, err := foo()
	if err != nil {
		return err
	}
	cmd.AddCommand(subCmd)
	return nil
}

func AddConfigCommand(cmd *cobra.Command, config *AppConfig, foo func(config *AppConfig) (*cobra.Command, error)) error {
	if config == nil {
		return fmt.Errorf(`app config is nil`)
	}
	subCmd, err := foo(config)
	if err != nil {
		return err
	}
	cmd.AddCommand(subCmd)
	return nil
}

const (
	HelpTmpl = `{{with (or .Long .Short)}}Name: {{. | trimTrailingWhitespaces}}

{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
	UsageTmpl = `Usage: {{if .Runnable}}{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}{{.CommandPath}} [OPTIONS] COMMAND{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Options:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Options:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} COMMAND --help" for more information about a command.{{end}}

To get more help with devtool, check out our guides at https://github.com/anthony-dong/golang
`
)
