package codec

import (
	"net/url"
	"strings"
	"unicode"

	"github.com/anthony-dong/golang/pkg/utils"

	"github.com/iancoleman/orderedmap"

	"github.com/anthony-dong/golang/pkg/codec"

	"github.com/mattn/go-shellwords"
)

func ParseUrl(_src []byte) (dst []byte, err error) {
	input := strings.TrimSpace(string(_src))
	if decode, err := codec.NewUrlCodec().Decode(utils.String2Bytes(input)); err == nil {
		input = string(decode)
	}
	parse, err := url.Parse(input)
	if err != nil {
		return nil, err
	}
	type URL struct {
		Scheme string                 `json:"Scheme"`
		User   *url.Userinfo          `json:"User"`  // username and password information
		Host   string                 `json:"Host"`  // host or host:port
		Path   string                 `json:"Path"`  // path (relative paths may omit leading slash)
		Query  *orderedmap.OrderedMap `json:"Query"` // encoded query values, without '?'
	}
	r := URL{
		Scheme: parse.Scheme,
		User:   parse.User,
		Host:   parse.Host,
		Path:   parse.Path,
		Query: func() *orderedmap.OrderedMap {
			result := orderedmap.New()
			query := parse.Query()
			for k, v := range query {
				switch len(v) {
				case 0:
					result.Set(k, nil)
				case 1:
					result.Set(k, v[0])
				default:
					result.Set(k, v)
				}
			}
			return result
		}(),
	}
	return []byte(utils.ToJson(r, true)), nil
}

func ParseCmd(src []byte) (dst []byte, err error) {
	cmdSlice, err := shellwords.Parse(string(src))
	if err == nil {
		fix := make([]string, 0, len(cmdSlice))
		for _, elem := range cmdSlice {
			// 会trim 全部的space， 如果data部分前缀有space的话确实会存在问题，但是对于命令行来说大部分不会影响
			fix = append(fix, strings.TrimLeftFunc(elem, func(r rune) bool {
				return unicode.IsSpace(r)
			}))
		}
		cmdSlice = fix
	}
	if len(cmdSlice) == 0 {
		return utils.ToJsonByte(cmdSlice), nil
	}
	switch cmdSlice[0] {
	case "curl":

		return nil, nil
	default:
		return utils.ToJsonByte(cmdSlice, true), nil
	}
}
