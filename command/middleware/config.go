package middleware

import (
	"encoding/json"
	"io/ioutil"
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
			if err := readConfig(config.ConfigFile, &config.Config); err != nil {
				return err
			}
			logs.Debug("init config success. filename: %s", config.ConfigFile)
			return nil
		}
		files := make([]string, 0)
		files = append(files, filepath.Join(utils.GetUserHomeDir(), command.UserHomeConfig))
		files = append(files, filepath.Join(utils.GetPwd(), command.CurrentDirConfig))
		for _, file := range files {
			if err := readConfig(file, &config.Config); err == nil {
				logs.Debug("init config success. filename: %s", file)
				return nil
			}
			continue
		}
		return nil
	}
}

func readConfig(filename string, cfg *command.Config) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(content, cfg); err == nil {
		return nil
	}
	return json.Unmarshal(content, cfg)
}
