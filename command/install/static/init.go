package static

import "embed"

//go:embed perf.sh
//go:embed sync.sh
//go:embed .fileignore
var embedFs embed.FS

func ReadFile(filename string) ([]byte, error) {
	return embedFs.ReadFile(filename)
}

func GetFiles() []string {
	return []string{
		"perf.sh",
		"sync.sh",
	}
}

func GetExtraFiles(filename string) (map[string][]byte, error) {
	switch filename {
	case "sync.sh":
		file, err := ReadFile(".fileignore")
		if err != nil {
			return nil, err
		}
		return map[string][]byte{
			"sync-devbox": file,
		}, nil
	}
	return nil, nil
}

var fileUsage = map[string]string{
	"sync.sh": `注意：记得使用前配置好环境变量!
export BYTED_USERNAME=xxx && export BYTED_HOST_IP=x.x.x.x`,
}

func GetUsage(filename string) string {
	return fileUsage[filename]
}
