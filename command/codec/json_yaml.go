package codec

import (
	"io"
	"os"

	"github.com/anthony-dong/golang/pkg/utils"
	"github.com/spf13/cobra"
)

func NewJson2YamlCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "json2yaml",
		Short: `convert json to yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !utils.CheckStdInFromPiped() {
				return cmd.Help()
			}
			all, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			output, err := utils.JsonToYaml(all)
			if err != nil {
				return err
			}
			if _, err := os.Stdout.Write(output); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd, nil
}

func NewYaml2JsonCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "yaml2json",
		Short: `convert yaml to json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !utils.CheckStdInFromPiped() {
				return cmd.Help()
			}
			all, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}
			output, err := utils.YamlToJson(all)
			if err != nil {
				return err
			}
			if _, err := os.Stdout.Write(output); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd, nil
}
