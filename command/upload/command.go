package upload

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"github.com/anthony-dong/golang/command"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/anthony-dong/golang/pkg/codec"
	"github.com/anthony-dong/golang/pkg/logs"
)

type uploadCommand struct {
	OssConfigType  string `json:"type,omitempty"`
	File           string `json:"file,omitempty"`
	FileNameDecode string `json:"decode,omitempty"`
	DstFile        string `json:"dst_file,omitempty"`
}

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{Use: "upload", Short: `File upload tool`}
	var (
		cfg = &uploadCommand{}
	)
	cmd.Flags().StringVarP(&cfg.File, "file", "f", "", "Set the local path of upload file")
	cmd.Flags().StringVarP(&cfg.FileNameDecode, "decode", "d", "uuid", "Set the upload file name decode method (uuid|url|md5)")
	cmd.Flags().StringVarP(&cfg.OssConfigType, "type", "t", "default", "Set the upload config type")
	cmd.Flags().StringVar(&cfg.DstFile, "dst", "", "Set the dst file name")
	if err := cmd.MarkFlagRequired("file"); err != nil {
		return nil, err
	}
	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		return cfg.Run(cmd.Context(), command.GetAppConfig(cmd.Context()).UploadConfig)
	}
	return cmd, nil
}

func (c *uploadCommand) validate() error {
	if c.File == "" {
		return fmt.Errorf("flag needs an argument: --file")
	}
	if filename, err := filepath.Abs(c.File); err != nil {
		return err
	} else {
		c.File = filename
	}
	if c.OssConfigType == "" {
		c.OssConfigType = "default"
	}
	logs.Info("[upload] config:\n%s", utils.ToJson(c, true))
	return nil
}
func (c *uploadCommand) Run(ctx context.Context, config *command.UploadConfig) error {
	if err := c.validate(); err != nil {
		return err
	}
	if config == nil || len(config.Bucket) == 0 {
		return errors.Errorf("not found bucket config, bucket: %s", c.OssConfigType)
	}
	cfg, isExist := config.Bucket[c.OssConfigType]
	if !isExist {
		return errors.Errorf(`invalid bucket type, type: %s`, c.OssConfigType)
	}
	_, suffix := utils.GetFilePrefixAndSuffix(c.File)
	prefix, err := c.getFileName(c.File)
	if err != nil {
		return errors.Errorf(`new file name err, err: %v`, err)
	}
	fileInfo := OssUploadFile{
		LocalFile:  c.File,
		SaveDir:    time.Now().Format(utils.FormatTimeV2),
		FilePrefix: prefix,
		FileSuffix: suffix,
		DstFile:    c.DstFile,
	}
	bucket, err := NewBucket(&cfg)
	if err != nil {
		return err
	}
	if err := fileInfo.PutFile(bucket, &cfg); err != nil {
		return err
	}
	fileUrl := fileInfo.GetOSSUrl(&cfg)
	if logs.IsLevel(logs.LevelInfo) {
		logs.Info("[upload] end success, url: %s", fileUrl)
	} else {
		fmt.Println(fileUrl)
	}
	return nil
}

func (c *uploadCommand) getFileName(filename string) (string, error) {
	switch c.FileNameDecode {
	case "uuid":
		return utils.GenerateUUID(), nil
	case "url":
		prefix, _ := utils.GetFilePrefixAndSuffix(filename)
		return string(codec.NewUrlCodec().Encode([]byte(prefix))), nil
	case "md5":
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			return "", fmt.Errorf(`read file: %s find err: %v`, filename, err)
		}
		return string(codec.NewMd5Codec().Encode(content)), nil
	}
	return utils.GenerateUUID(), nil
}
