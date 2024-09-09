package _init

import (
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/utils"
)

func NewVscodeLaunch() (*cobra.Command, error) {
	output := ""
	force := false
	cmd := &cobra.Command{
		Use:   "vlgo",
		Short: `the go vscode launch template`,
		RunE: func(cmd *cobra.Command, args []string) error {
			files := loadGoLaunchFile()
			for _, file := range files {
				if err := file.Write(output, force); err != nil {
					return err
				}
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&output, "output", "O", utils.GetPwd(), "the output dir")
	cmd.Flags().BoolVarP(&force, "force", "f", false, "force replace file")
	return cmd, nil
}

func loadGoLaunchFile() []*File {
	return []*File{
		{
			Name:       ".vscode/launch.json",
			IsTemplate: false,
			Content: `{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "example_vscode",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}/example/vscode/model.go",
            "args": [
                ""
            ],
            "cwd": "${workspaceFolder}",
            "env": {
                "K1": "V1"
            }
        },
        {
            "name": "example_vscode_test",
            "type": "go",
            "request": "launch",
            "mode": "test",
            "program": "${workspaceFolder}/example/vscode",
            "args": [
                "-test.run=Test_main"
            ],
            "cwd": "${workspaceFolder}",
            "env": {
                "K1": "V1"
            }
        }
    ]
}`,
		},
	}
}
