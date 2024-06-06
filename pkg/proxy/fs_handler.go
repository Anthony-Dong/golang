package proxy

import (
	"bytes"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/fatih/color"

	"github.com/anthony-dong/golang/pkg/logs"
	"github.com/anthony-dong/golang/pkg/utils"
)

const textContentType = "text/plain; charset=utf-8"
const contentTypeHeader = "Content-Type"

type FsHandler struct {
	dir         string
	handler     http.Handler
	dirTemplate *template.Template
}

func NewFsHandler(dir string) http.Handler {
	f := FsHandler{dir: dir}
	f.handler = http.FileServer(http.Dir(f.dir))
	f.dirTemplate = utils.MustTemplate("", map[string]interface{}{
		"FormatSize": func(size int) string {
			if size == 0 {
				return "-"
			}
			return utils.FormatSize(size)
		},
		"FormatTime": func(t time.Time) string {
			if t.Unix() == 0 {
				return "-"
			}
			return t.Format("2006-01-02 15:04:05")
		},
		"Quote": func(s string) string {
			return strconv.Quote(s)
		},
	}, `<!DOCTYPE html>
<html lang="zh-CN">
<head>
	<meta charset="UTF-8">
    <title>目录信息</title>
    <style>
        table {
            width: 100%;
            border-collapse: collapse;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        tr:nth-child(even) {
            background-color: #f2f2f2;
        }
        th {
            background-color: #4CAF50;
            color: white;
        }
		a {
            text-decoration: none;
        }
        /* 这将改变所有链接的颜色为黑色，并去掉下划线 */
        a:link, a:visited {
            color: black;
            text-decoration: none;
        }
        /* 鼠标悬停在链接上时，将颜色改变为蓝色，并添加下划线 */
        a:hover, a:active {
            color: blue;
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <h2>目录信息 Total: {{ .Total }} </h2>
    <table>
        <tr>
            <th>文件名</th>
            <th>文件大小</th>
            <th>最后修改时间</th>
        </tr>
{{- range .Files }}
        <tr>
			<td><a href={{ Quote .LinkName }}>{{ .LinkName }}</a></td>
            <td>{{ FormatSize .Size }}</td>
            <td>{{ FormatTime .ModTime }}</td>
        </tr>
{{- end }}
    </table>
</body>
</html>`)
	return &f
}

func (f *FsHandler) GetAbsPath(name string) (string, error) {
	if name == "" || name == "/index.html" {
		name = "/"
	}
	name = filepath.Clean(name)
	if filepath.Separator != '/' && strings.ContainsRune(name, filepath.Separator) {
		return "", errors.New("http: invalid character in file path")
	}
	dir := f.dir
	if dir == "" {
		dir = "."
	}
	return filepath.Join(dir, filepath.FromSlash(path.Clean("/"+name))), nil
}

func (f *FsHandler) HandleDirectory(dir string, w http.ResponseWriter) error {
	entries, err := utils.ReadDir(dir)
	if err != nil {
		return err
	}
	buffer := bytes.Buffer{}
	fileInfos := make([]map[string]interface{}, 0, len(entries))
	fileInfos = append(fileInfos, map[string]interface{}{
		"Name":     "../",
		"LinkName": "../",
		"Size":     0,
		"ModTime":  time.Unix(0, 0),
		"IsDir":    true,
	})
	for _, elem := range entries {
		fileInfo, err := elem.Info()
		if err != nil {
			return err
		}
		fileName := filepath.Base(fileInfo.Name())
		linkName := fileName
		fileSize := fileInfo.Size()
		if fileInfo.IsDir() {
			linkName = linkName + "/"
			fileSize = 0
		}
		fileInfos = append(fileInfos, map[string]interface{}{
			"Name":     fileName,
			"LinkName": linkName,
			"Size":     int(fileSize),
			"ModTime":  fileInfo.ModTime(),
			"IsDir":    fileInfo.IsDir(),
		})
	}
	if err := f.dirTemplate.Execute(&buffer, map[string]interface{}{
		"Files": fileInfos,
		"Total": len(fileInfos),
	}); err != nil {
		return err
	}
	w.Header().Set(contentTypeHeader, "text/html; charset=utf-8")
	_, _ = w.Write(buffer.Bytes())
	return nil
}

func (f *FsHandler) GetFileContentType(fileInfo fs.File) string {
	buffer := make([]byte, 512)
	if readSize, err := fileInfo.Read(buffer); err != nil {
		return ""
	} else {
		buffer = buffer[:readSize]
	}
	return http.DetectContentType(buffer)
}

func (f *FsHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if err := f.customHandle(writer, request); err != nil {
		f.handler.ServeHTTP(writer, request)
	}
	logs.CtxInfo(nil, "%s %s %s", request.Method, color.CyanString(request.URL.Path), color.WhiteString(formatContentType(writer.Header().Get(contentTypeHeader))))
}

func formatContentType(contentType string) string {
	split := strings.Split(contentType, ";")
	return strings.TrimSpace(split[0])
}

func (f *FsHandler) customHandle(w http.ResponseWriter, r *http.Request) error {
	var downGradeError = errors.New(`downGradeError`)
	absPath, err := f.GetAbsPath(r.URL.Path)
	if err != nil {
		return downGradeError
	}
	fileInfo, err := os.Open(absPath)
	if err != nil {
		return downGradeError
	}
	defer fileInfo.Close()
	stat, err := fileInfo.Stat()
	if err != nil {
		return downGradeError
	}
	if stat.IsDir() {
		return f.HandleDirectory(absPath, w)
	}
	if c := f.GetFileContentType(fileInfo); c != "" {
		w.Header().Set(contentTypeHeader, c)
	}
	return downGradeError
}
