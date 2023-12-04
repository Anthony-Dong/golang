package git

import (
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "git",
		Short: "The git tools",
	}
	if err := command.AddCommand(cmd, NewGitCloneCommand); err != nil {
		return nil, err
	}
	return cmd, nil
}
