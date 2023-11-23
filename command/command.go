package command

import (
	"github.com/spf13/cobra"
)

type Middleware func(cmd *cobra.Command, args []string) error
