package tcpdump

import (
	"encoding/hex"

	"github.com/anthony-dong/golang/pkg/utils"
)

func s2b(s string) []byte {
	return utils.String2Bytes(s)
}

func b2s(b []byte) string {
	return utils.Bytes2String(b)
}

func HexDump(src []byte) string {
	return hex.Dump(src)
}
