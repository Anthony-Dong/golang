package _goto

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "goto",
		Short: "goto anywhere",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(ParseGoLink(args[0]))
			return nil
		},
	}
	return cmd, nil
}

func ParseGoLink(input string) string {
	var urlRegexp = regexp.MustCompile(`([\w_.-]+/[\w_.-]+/[\w_.-]+)/?([\w_.@-]+)/(.*):(\d+)`)

	submatch := urlRegexp.FindStringSubmatch(input)
	if len(submatch) != 5 {
		return ""
	}
	repo := submatch[1]
	tag := submatch[2]
	file := submatch[3]
	line := submatch[4]

	//fmt.Println("repo:", repo, "tag:", tag, "file:", file, "line:", line)
	return fmt.Sprintf("https://%s/blob/%s/%s#L%s", repo, replaceUrlTag(tag), file, line)
}

func replaceUrlTag(tag string) string {
	if strings.HasPrefix(tag, "@") {
		tag = strings.TrimPrefix(tag, "@")
	}
	urlTag := strings.ReplaceAll(tag, "@", "/")
	urlTagSplit := strings.Split(urlTag, "/")
	tagPrefix := strings.Join(urlTagSplit[:len(urlTagSplit)-1], "/")
	if tagPrefix != "" {
		urlTag = urlTag + "/" + tagPrefix
	}
	return urlTag
}
