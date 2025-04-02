package http2

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
	"testing"
)

func TestDecodeH2CMessage(t *testing.T) {
	data := `505249202a20485454502f322e300d0a0d0a534d0d0a0d0a
00000604000000000000040000ffff
0000520104000000018386459060a49cebef2a13b8b677310ac6024e7f4192416cee5b22a87496294f645219043f72d9eb5f911d75d0620d263d4c4d6564ff699ec34a9f7a8aea6497cb1dc0b844b83f40027465864d833505b11f000000040100000000
00001600000000000100000000110b000100000009636c69656e742c20300000001600000000000100000000110b000100000009636c69656e742c20310000001600000000000100000000110b000100000009636c69656e742c203200
000000000100000001
0000040800000000000000001600000806000000000002041010090e070700000806010000000002041010090e0707`

	hexData, err := parseHexData(data)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(base64.StdEncoding.EncodeToString(hexData))
	r := bytes.NewReader(hexData)
	reader := bufio.NewReader(r)
	message, err := DecodeH2CMessage(reader)
	if err != nil {
		if err == io.EOF {
			t.Fatal(err)
		}
	}
	t.Log(message)

	for {

		frame, err := DecodeFrame(reader)
		if err != nil {
			if err == io.EOF {
				break
			}
			t.Fatal(err)
		}
		t.Log(ConvertToFrame(frame))
	}

	t.Log("end")
}

func parseHexData(hexData string) ([]byte, error) {
	// 去除换行符
	hexData = strings.ReplaceAll(hexData, "\n", "")
	// 将十六进制字符串转换为字节切片
	bytes, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
