package cpp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_readFileArgs(t *testing.T) {
	args, err := readBuildAndLinkArgs(bytes.NewBufferString(`
// link: -lspdlog
// build: -O2
// cxxopt: -g
// link: -L/usr/local/lib
// linkopt: -lgtest
`))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v\n", args)
	assert.Equal(t, args, &buildAndLink{buildArgs: []string{"-O2", "-g"}, linkArgs: []string{"-lspdlog", "-L/usr/local/lib", "-lgtest"}})
}
