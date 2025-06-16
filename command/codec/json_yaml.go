package codec

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/utils"
)

func NewJson2YamlCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "json2yaml",
		Short: "Convert JSON data to YAML format",
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
		Short: "Convert YAML data to JSON format",
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
