package command

import (
	"path/filepath"

	"github.com/anthony-dong/golang/pkg/utils"
)

const AppConfigFile = "devtool.yaml"
const AppHomeDir = ".devtool"

const AppVersion = "v0.0.3"
const AppName = "devtool"

func GetAppHomeDir() string {
	return filepath.Join(utils.GetUserHomeDir(), AppHomeDir)
}
