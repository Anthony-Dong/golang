package lint

import (
	"bytes"
	"fmt"
	"go/scanner"
	"go/token"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestFormatYaml(t *testing.T) {
	output, err := FormatYaml([]byte(`
---
k1: v1 # 111
k2: v2
...

--- # 最喜愛的電影
k1: v11
k2: v22
...
`))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, string(output), `k1: v1 # 111
k2: v2
---
# 最喜愛的電影
k1: v11
k2: v22
`)
}

func TestName(t *testing.T) {
	// src is the input that we want to tokenize.     src := [] byte("cos(x) + 1i*sin(x) // Euler")      // Initialize the scanner.     var s scanner. Scanner     fset := token. NewFileSet()                      // positions are relative to fset     file := fset. AddFile("", fset. Base(), len(src)) // register input "file"     s. Init(file, src, nil /* no error handler */, scanner. ScanComments)      // Repeated calls to Scan yield the token sequence found in the input.     for {         pos, tok, lit := s. Scan()         if tok == token. EOF {             break         }         fmt. Printf("%s\t%s\t%q\n", fset. Position(pos), tok, lit)     }
	// src is the input that we want to tokenize.
	src := []byte(`
# 框架配置
struct BizFlexConfig {
    1: string psm
    2: string desc
    3: list<string> owners
    4: list<PluginInfo> plugins
    5: list<SceneInfo> scenes
}`)

	// Initialize the scanner.
	var s scanner.Scanner
	fset := token.NewFileSet()                      // positions are relative to fset
	file := fset.AddFile("", fset.Base(), len(src)) // register input "file"
	s.Init(file, src, nil /* no error handler */, scanner.ScanComments)

	// Repeated calls to Scan yield the token sequence found in the input.
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		fmt.Printf("%s %s %q\n", fset.Position(pos), tok, lit)
	}

}

type formater struct {
	lines []bytes.Buffer
}

func (f *formater) write(pos token.Pos, tok token.Token, lit string) {
	switch tok {
	case token.ILLEGAL:
		if lit == "#" {

		}
	default:
		return
	}
}

func FormatYaml(input []byte) ([]byte, error) {
	if len(input) == 0 {
		return []byte{}, nil
	}
	decoder := yaml.NewDecoder(bytes.NewBuffer(input))
	output := bytes.NewBuffer(make([]byte, 0, len(input)))
	encoder := yaml.NewEncoder(output)
	encoder.SetIndent(4)
	for {
		node := yaml.Node{}
		if err := decoder.Decode(&node); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if err := encoder.Encode(&node); err != nil {
			return nil, err
		}
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}
