package proxy

import (
	"bytes"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sort"
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
            <th>文件名<span onclick="window.location.href = updateURLParameter(window.location.href, 'sort', 'name')"> ⬇ </span></th>
            <th>文件大小<span onclick="window.location.href = updateURLParameter(window.location.href, 'sort', 'size')"> ⬇ </span></th>
            <th>最后修改时间<span onclick="window.location.href = updateURLParameter(window.location.href, 'sort', 'mtime')"> ⬇ </span></th>
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
<script>
    function updateURLParameter(url, param, paramVal) {
        var newAdditionalURL = "";
        var tempArray = url.split("?");
        var baseURL = tempArray[0];
        var additionalURL = tempArray[1] ? tempArray[1].split("&") : [];

        var temp = "";
        for (var i = 0; i < additionalURL.length; i++) {
            if (additionalURL[i].split('=')[0] != param) {
                newAdditionalURL += temp + additionalURL[i];
                temp = "&";
            }
        }

        var rows_txt = temp + "" + param + "=" + paramVal;
        return baseURL + "?" + newAdditionalURL + rows_txt;
    }
</script>
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

type FileInfo struct {
	Name     string
	LinkName string
	Size     int
	ModTime  time.Time
	IsDir    bool
}

func (f *FsHandler) HandleDirectory(dir string, w http.ResponseWriter, sortKey string) error {
	entries, err := utils.ReadDir(dir)
	if err != nil {
		return err
	}
	buffer := bytes.Buffer{}
	fileInfos := make([]*FileInfo, 0, len(entries))
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
		fileInfos = append(fileInfos, &FileInfo{
			Name:     fileName,
			LinkName: linkName,
			Size:     int(fileSize),
			ModTime:  fileInfo.ModTime(),
			IsDir:    fileInfo.IsDir(),
		})
	}

	switch sortKey {
	case "mtime":
		sort.Slice(fileInfos, func(i, j int) bool {
			return fileInfos[i].ModTime.After(fileInfos[j].ModTime)
		})
	case "size":
		sort.Slice(fileInfos, func(i, j int) bool {
			if fileInfos[i].Size == fileInfos[j].Size {
				return fileInfos[i].Name < fileInfos[j].Name
			}
			return fileInfos[i].Size < fileInfos[j].Size
		})
	default:
		sort.Slice(fileInfos, func(i, j int) bool {
			return strings.ToLower(fileInfos[i].Name) < strings.ToLower(fileInfos[j].Name)
		})
	}
	fileInfos = append([]*FileInfo{
		{
			Name:     "../",
			LinkName: "../",
			Size:     0,
			ModTime:  time.Unix(0, 0),
			IsDir:    true,
		},
	}, fileInfos...)

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
	if stat, _ := fileInfo.Stat(); stat != nil {
		switch filepath.Ext(stat.Name()) {
		case ".json":
			return "application/json"
		case ".svg":
			return "image/svg+xml"
		case ".xml":
			return textContentType // xml 浏览器渲染有问题，这里当作文本处理
		}
	}
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
	sortKey := r.URL.Query().Get("sort")
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
		return f.HandleDirectory(absPath, w, sortKey)
	}
	if c := f.GetFileContentType(fileInfo); c != "" {
		w.Header().Set(contentTypeHeader, c)
	}
	return downGradeError
}
