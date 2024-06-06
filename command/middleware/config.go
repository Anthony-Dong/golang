package middleware

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewInitConfigMv(configFile string, config *command.AppConfig) command.Middleware {
	return func(cmd *cobra.Command, args []string) error {
		if configFile != "" {
			if err := utils.UnmarshalFromFile(configFile, &config.CommandConfig); err != nil {
				return err
			}
			logs.Debug("init config success. filename: %s", configFile)
			return nil
		}
		files := make([]string, 0)
		executable, err := os.Executable()
		if err != nil {
			return err
		}
		files = append(files, filepath.Join(utils.GetPwd(), command.AppConfigFile))
		files = append(files, filepath.Join(filepath.Dir(executable), command.AppConfigFile))
		files = append(files, filepath.Join(command.GetAppHomeDir(), command.AppConfigFile))
		for _, file := range files {
			if err := utils.UnmarshalFromFile(file, &config.CommandConfig); err == nil {
				logs.Debug("init config success. filename: %s", file)
				return nil
			}
			continue
		}
		return nil
	}
}
