package command

import (
	"context"

	"github.com/anthony-dong/golang/pkg/rpc"
)

type AppConfig struct {
	Verbose    bool
	LogLevel   string
	ConfigFile string
	AppName    string
	AppVersion string

	StaticConfig
	CurlConfig *CurlConfig

	//Middlewares []func(config *AppConfig) Middleware
}

type StaticConfig struct {
	UploadConfig  *UploadConfig  `yaml:"Upload,omitempty"`
	HexoConfig    *HexoConfig    `yaml:"Hexo,omitempty"`
	RunTaskConfig *RunTaskConfig `yaml:"RunTask,omitempty"`

	Middlewares []func(config *AppConfig) Middleware `yaml:"-"`
}

type RunTaskConfig struct {
	Includes []string `yaml:"Includes,omitempty"`
}

type UploadConfig struct {
	Bucket map[string]OSSConfig `yaml:"Bucket,omitempty"`
}

type OSSConfig struct {
	AccessKeyId     string `yaml:"AccessKeyId,omitempty"`
	AccessKeySecret string `yaml:"AccessKeySecret,omitempty"`
	Endpoint        string `yaml:"Endpoint,omitempty"`
	UrlEndpoint     string `yaml:"UrlEndpoint,omitempty"`
	Bucket          string `yaml:"Bucket,omitempty"`
	PathPrefix      string `yaml:"PathPrefix,omitempty"`
}

type HexoConfig struct {
	Ignore  []string `yaml:"Ignore,omitempty"`
	KeyWord []string `yaml:"KeyWord,omitempty"`
}

type CurlConfig struct {
	NewClient func(ctx context.Context, request *rpc.Request, idl *rpc.IDLInfo) (rpc.Client, error)
}
