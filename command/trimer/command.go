package trimer

import (
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "trimer",
		Short: "Trim or filter data from various formats",
	}
	if err := command.AddCommand(cmd, NewJsonTrimerCommand); err != nil {
		return nil, err
	}
	return cmd, nil
}
