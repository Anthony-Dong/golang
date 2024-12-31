package cpp

import (
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/anthony-dong/golang/pkg/utils"
)

type buildAndLink struct {
	buildArgs []string
	linkArgs  []string
}

var buildArgsRegexp = regexp.MustCompile(`^//\s*(build|cxxopt)\s*:\s*`)
var linkArgsRegexp = regexp.MustCompile(`^//\s*(link|linkopt)\s*:\s*`)

func readFileArgs(file string) (*buildAndLink, error) {
	open, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer open.Close()
	return readBuildAndLinkArgs(open)
}

func readBuildAndLinkArgs(reader io.Reader) (*buildAndLink, error) {
	result := &buildAndLink{}
	utils.ReadSomeLines(reader, func(index int, line string) bool {
		if strings.HasPrefix(line, "#include") {
			return false
		}
		if prefix := linkArgsRegexp.FindString(line); prefix != "" {
			args := strings.TrimPrefix(line, prefix)
			result.linkArgs = append(result.linkArgs, utils.SplitArgs(args)...)
		}
		if prefix := buildArgsRegexp.FindString(line); prefix != "" {
			args := strings.TrimPrefix(line, prefix)
			result.buildArgs = append(result.buildArgs, utils.SplitArgs(args)...)
		}
		return true
	})
	return result, nil
}
