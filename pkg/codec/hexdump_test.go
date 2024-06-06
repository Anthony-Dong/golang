package codec

import (
	"encoding/base64"
	"strconv"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

var testcase = map[string]string{
	`00000000  0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae     (     0       `: `0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae`,
	`00019740  35 30 30 37 42 34 45 31 35 35                     5007B4E155      `: `35 30 30 37 42 34 45 31 35 35`,
	`00000030  63 92 01 0b 32 09 2a 07 23 46 30 46 30 46 30 a0   c   2 * #F0F0F0 `: `63 92 01 0b 32 09 2a 07 23 46 30 46 30 46 30 a0`,
	`	0x0040:  eafc b74a eafc b74a 4745 5420 2f68 656c  ...J...JGET./hel`: "eafc b74a eafc b74a 4745 5420 2f68 656c",
	`00000000  0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae     (     0`:                                                                                         "0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae",
	"0x0000:  600e 3d55 0020 0640 0000 0000 0000 0000  `.=U...@........":                                                                                            "600e 3d55 0020 0640 0000 0000 0000 0000",
	`00:02:30.058133 IP6 localhost.36962 > localhost.smc-https: Flags [P.], seq 1:84, ack 1, win 43, options [nop,nop,TS val 3942430538 ecr 3942430538], length 83`: "",
}

func TestReadHexdump2(t *testing.T) {
	testIsHex(t, `00000000  0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae     (     0       `, "0a 07 28 ef 9a ac de b1 30 12 bd ae 06 0a 93 ae", true)
}

func isEqual(str1, str2 string) bool {
	r1 := strings.Builder{}
	for _, elem := range str1 {
		if unicode.IsSpace(elem) {
			continue
		}
		r1.WriteRune(elem)
	}

	r2 := strings.Builder{}
	for _, elem := range str2 {
		if unicode.IsSpace(elem) {
			continue
		}
		r2.WriteRune(elem)
	}
	return r1.String() == r2.String()
}

func testIsHex(t testing.TB, k, v string, isCheck bool) {
	hexdump, isEnd := ReadHexdump(k)
	if !isCheck {
		return
	}
	if v == "" || hexdump == "" {
		assert.Equal(t, v, "")
		assert.Equal(t, hexdump, "")
		assert.Equal(t, isEnd, false)
		return
	}
	vs := strings.Builder{}
	for _, elem := range v {
		if unicode.IsSpace(elem) {
			continue
		}
		vs.WriteRune(elem)
	}
	assert.Equal(t, isEnd, len(vs.String()) < 32, k)
	assert.Equal(t, isEqual(hexdump, v), true, k)
}

func TestReadHexdump(t *testing.T) {
	for k, v := range testcase {
		testIsHex(t, k, v, true)
	}
}

func TestReadInt(t *testing.T) {
	t.Run("test1", func(t *testing.T) {
		parseInt, _ := strconv.ParseInt("0x0010", 0, 64)
		t.Log(parseInt)
	})
	// 00000010
	t.Run("test1", func(t *testing.T) {
		parseInt, _ := strconv.ParseInt("0x00000010", 0, 64)
		t.Log(parseInt)
	})
}

func BenchmarkReadHexdump(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for k, v := range testcase {
			testIsHex(b, k, v, false)
		}
	}
}

func Test_isByte(t *testing.T) {
	assert.Equal(t, isByte(1), true)
	assert.Equal(t, isByte(256), false)
}

func TestReadFile(t *testing.T) {
	dst, err := NewHexDumpCodec().Decode([]byte(`00000000  00 00 00 cc 10 00 00 00  00 00 00 00 00 01 00 00  |................|
00000010  00 00 00 00 00 ba 80 01  00 02 00 00 00 0c 70 72  |..............pr|
00000020  65 64 69 63 74 5f 6c 69  74 65 00 00 00 00 0c 00  |edict_lite......|
00000030  00 0f 00 01 0c 00 00 00  01 0a 00 01 00 06 5f 27  |.............._'|
00000040  be 98 10 32 0d 00 03 0b  13 00 00 00 02 00 00 00  |...2............|
00000050  09 61 64 73 5f 62 72 61  69 6e 3f 80 00 00 00 00  |.ads_brain?.....|
00000060  00 11 61 64 73 5f 62 72  61 69 6e 3a 61 64 5f 73  |..ads_brain:ad_s|
00000070  74 6f 70 3f 80 00 00 0d  00 04 0b 04 00 00 00 02  |top?............|
00000080  00 00 00 09 61 64 73 5f  62 72 61 69 6e 3f f0 00  |....ads_brain?..|
00000090  00 00 00 00 00 00 00 00  11 61 64 73 5f 62 72 61  |.........ads_bra|
000000a0  69 6e 3a 61 64 5f 73 74  6f 70 3f f0 00 00 00 00  |in:ad_stop?.....|
000000b0  00 00 00 0c 00 ff 0b 00  01 00 00 00 00 08 00 02  |................|
000000c0  00 00 00 00 0d 00 03 0b  0b 00 00 00 00 00 00 00  |................|`))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(base64.StdEncoding.EncodeToString(dst))
}
