package upload

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

func NewDownloadCommand(clientGetter func() (Client, error)) (*cobra.Command, error) {
	downloader := fileDownloader{}
	cmd := &cobra.Command{
		Use:   "download",
		Short: "download a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientGetter()
			if err != nil {
				return err
			}
			downloader.client = client
			return downloader.download(cmd.Context())
		},
	}
	cmd.Flags().StringVarP(&downloader.outputFile, "output", "o", "", "指定输出的文件名，默认object key的名称")
	cmd.Flags().IntVar(&downloader.bench, "bench", 10, "下载文件的并发数")
	cmd.Flags().StringVarP(&downloader.objectKey, "key", "k", "", "指定 object key")
	if err := cmd.MarkFlagRequired("key"); err != nil {
		return nil, err
	}
	return cmd, nil
}

type fileDownloader struct {
	objectKey  string
	outputFile string
	bench      int
	client     Client
}

func (f *fileDownloader) download(ctx context.Context) (err error) {
	if f.outputFile == "" {
		f.outputFile = f.objectKey
	}
	if abs, err := filepath.Abs(f.outputFile); err != nil {
		return err
	} else {
		f.outputFile = abs
	}
	if utils.ExistFile(f.outputFile) {
		return fmt.Errorf("file %q already exists", f.outputFile)
	}
	logs.CtxInfo(ctx, "download object %q into %q", f.objectKey, f.outputFile)

	start := time.Now()
	defer func() {
		if err == nil {
			logs.CtxInfo(ctx, "downloading object %s success, spend: %s", f.objectKey, time.Since(start))
		}
	}()

	files, err := f.downloadSingleFile(ctx)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return nil
	}
	return f.downloadMultiFile(ctx, files)
}

func (f *fileDownloader) downloadSingleFile(ctx context.Context) ([]ObjectInfo, error) {
	file, err := os.OpenFile(f.outputFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer file.Close()
	if err := f.client.GetObject(ctx, file, f.objectKey); err != nil {
		return nil, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	objects := make([]ObjectInfo, 0)
	if err := json.NewDecoder(file).Decode(&objects); err != nil {
		return nil, nil
	}
	return objects, nil
}

func (f *fileDownloader) downloadMultiFile(ctx context.Context, objects []ObjectInfo) error {
	multiFiles := make([]*utils.SplitFile, 0)
	for _, object := range objects {
		multiFiles = append(multiFiles, utils.NewSplitFile(f.outputFile, object.Start, object.End))
	}
	wg := errgroup.Group{}
	if f.bench > 1 {
		wg.SetLimit(f.bench)
	} else {
		wg.SetLimit(1)
	}
	for index, _ := range multiFiles {
		index := index
		wg.Go(func() error {
			start := time.Now()

			file := multiFiles[index]
			if err := file.Init(); err != nil {
				return err
			}
			defer file.Close()
			objectKey := objects[index].Key
			if err := f.client.GetObject(ctx, file, objectKey); err != nil {
				return errors.Wrapf(err, "download object %q error", objectKey)
			}
			logs.CtxInfo(ctx, "download object success. object: %s, spend: %s", objectKey, time.Since(start))
			return nil
		})
	}
	return wg.Wait()
}

func NewUploadCommand(clientGetter func() (Client, error)) (*cobra.Command, error) {
	uploader := fileUploader{}
	cmd := &cobra.Command{
		Use:   "upload [-f file]",
		Short: "upload a file",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := clientGetter()
			if err != nil {
				return err
			}
			uploader.client = client
			return uploader.upload(cmd.Context())
		},
	}
	cmd.Flags().StringVarP(&uploader.objectKey, "key", "k", "", "指定文件的object key, 默认是文件名")
	cmd.Flags().StringVarP(&uploader.inputFile, "file", "f", "", "需要上传的文件")
	cmd.Flags().IntVar(&uploader.split, "split", 1024*1024*50, "文件分片大小，默认50MB")
	cmd.Flags().IntVar(&uploader.bench, "bench", 10, "上传文件的并发数")
	if err := cmd.MarkFlagRequired("file"); err != nil {
		return nil, err
	}
	return cmd, nil
}

type fileUploader struct {
	client Client

	inputFile string

	objectKey string

	bench int
	split int
}

func (f *fileUploader) upload(ctx context.Context) error {
	if abs, err := filepath.Abs(f.inputFile); err != nil {
		return err
	} else {
		f.inputFile = abs
	}
	if f.objectKey == "" {
		f.objectKey = filepath.Base(f.inputFile)
	}
	logs.CtxInfo(ctx, "uploading file %q into object %q", f.inputFile, f.objectKey)
	files, err := utils.NewMultiSplitFile(f.inputFile, f.split)
	if err != nil {
		return errors.Wrapf(err, "split file %s error", f.inputFile)
	}
	if len(files) > 1 {
		return f.uploadMultiObject(ctx, files)
	}
	return f.uploadSingleObject(ctx)
}

func (f *fileUploader) uploadSingleObject(ctx context.Context) error {
	start := time.Now()
	file, err := os.Open(f.inputFile)
	if err != nil {
		return errors.Wrapf(err, "open file %s error", f.inputFile)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return errors.Wrapf(err, "stat file %s error", f.inputFile)
	}
	if err := f.client.PutObject(ctx, file, f.objectKey, int(stat.Size())); err != nil {
		return errors.Wrapf(err, "put object %s error", f.objectKey)
	}
	logs.CtxInfo(ctx, "put object success. object: %s, size: %v, spend: %s", f.objectKey, stat.Size(), time.Since(start))
	return nil
}

func (f *fileUploader) uploadMultiObject(ctx context.Context, files []*utils.SplitFile) error {
	totalStart := time.Now()
	objects := make([]ObjectInfo, 0)
	wg := errgroup.Group{}
	if f.bench > 1 {
		wg.SetLimit(f.bench)
	} else {
		wg.SetLimit(1)
	}
	for index, file := range files {
		objectKey := f.objectKey + "." + strconv.Itoa(index)
		object := ObjectInfo{Key: objectKey}
		object.Start, object.End = file.Index()
		objects = append(objects, object)

		file := file
		wg.Go(func() error {
			start := time.Now()
			if err := file.Init(); err != nil {
				return errors.Wrapf(err, "init file %s error", file.FileName())
			}
			defer file.Close()
			if err := f.client.PutObject(ctx, file, objectKey, int(file.Size())); err != nil {
				return errors.Wrapf(err, "put object %s error", objectKey)
			}
			logs.CtxInfo(ctx, "put object success. object: %s, size: %d, spend: %s", objectKey, file.Size(), time.Since(start))
			return nil
		})
	}
	if err := wg.Wait(); err != nil {
		return err
	}
	start := time.Now()
	marshal, err := json.Marshal(objects)
	if err != nil {
		return err
	}
	data := bytes.NewBuffer(marshal)
	if err := f.client.PutObject(ctx, data, f.objectKey, data.Len()); err != nil {
		return errors.Wrapf(err, "put object %s error", f.objectKey)
	}
	logs.CtxInfo(ctx, "put object success. object: %s, spend: %s, total_spend: %s", f.objectKey, time.Since(start), time.Since(totalStart))
	return nil
}

type ObjectInfo struct {
	Key   string `json:"object_key"`
	Start int    `json:"start"`
	End   int    `json:"end"`
}

func (f ObjectInfo) String() string {
	marshal, _ := json.Marshal(f)
	return string(marshal)
}

type Client interface {
	PutObject(ctx context.Context, r io.Reader, key string, size int) error
	GetObject(ctx context.Context, w io.Writer, key string) error
}
