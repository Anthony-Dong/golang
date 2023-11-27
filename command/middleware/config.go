package middleware

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/anthony-dong/golang/command"
	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

func NewInitConfigMv(config *command.AppConfig) Middleware {
	return func(cmd *cobra.Command, args []string) error {
		if config.ConfigFile != "" {
			if err := readConfig(config.ConfigFile, config); err != nil {
				return err
			}
			logs.Debug("init config success. filename: %s", config.ConfigFile)
			return nil
		}
		files := make([]string, 0)
		executable, err := os.Executable()
		if err != nil {
			return err
		}
		files = append(files, filepath.Join(utils.GetPwd(), command.AppConfigFile))
		files = append(files, filepath.Join(utils.GetUserHomeDir(), command.UserHomeConfigFile))
		files = append(files, filepath.Join(filepath.Dir(executable), command.AppConfigFile))
		for _, file := range files {
			if err := readConfig(file, config); err == nil {
				logs.Debug("init config success. filename: %s", file)
				return nil
			}
			continue
		}
		return nil
	}
}

func readConfig(filename string, cfg *command.AppConfig) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, cfg)
}
