package command

type AppConfig struct {
	Verbose    bool
	LogLevel   string
	ConfigFile string
	Config
}

type Config struct {
	UploadConfig  UploadConfig  `yaml:"Upload"`
	HexoConfig    HexoConfig    `yaml:"Hexo"`
	RunTaskConfig RunTaskConfig `yaml:"RunTask"`
}

type RunTaskConfig struct {
	Includes []string `yaml:"Includes"`
}

type UploadConfig struct {
	Bucket map[string]OSSConfig `yaml:"Bucket"`
}

type OSSConfig struct {
	AccessKeyId     string `yaml:"AccessKeyId"`
	AccessKeySecret string `yaml:"AccessKeySecret"`
	Endpoint        string `yaml:"Endpoint"`
	UrlEndpoint     string `yaml:"UrlEndpoint"`
	Bucket          string `yaml:"Bucket"`
	PathPrefix      string `yaml:"PathPrefix"`
}

type HexoConfig struct {
	Ignore  []string `yaml:"Ignore,omitempty"`
	KeyWord []string `yaml:"KeyWord,omitempty"`
}
