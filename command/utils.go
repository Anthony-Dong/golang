package command

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/logs"
)

func ExitError(err error) {
	if err != nil {
		logs.Error("[%s] exit error: %s", AppName, err.Error())
	}
	os.Exit(1)
}

func AddCommandWithConfig[T any](cmd *cobra.Command, config func() *T, foo func(config func() *T) (*cobra.Command, error)) error {
	subCmd, err := foo(config)
	if err != nil {
		return err
	}
	cmd.AddCommand(subCmd)
	return nil
}

func AddCommand(cmd *cobra.Command, foo func() (*cobra.Command, error)) error {
	subCmd, err := foo()
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
  {{.Example | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableSubCommands}}

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
