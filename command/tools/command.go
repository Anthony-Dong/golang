package tools

import (
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{Use: "tools", Short: "misc tools"}
	if err := command.AddCommand(cmd, NewClangdCommand); err != nil {
		return nil, err
	}
	return cmd, nil
}
