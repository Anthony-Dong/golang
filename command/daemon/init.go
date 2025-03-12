package bstatic

import (
	"embed"
)

//go:embed mac_daemon.plist
//go:embed linux_daemon.service
var fs embed.FS

const LinuxDaemonFile = "linux_daemon.service"
const MacDaemonFile = "mac_daemon.plist"

func ReadFile(name string) (string, error) {
	file, err := fs.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
