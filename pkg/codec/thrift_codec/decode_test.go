package thrift_codec

import (
	"bytes"
	"context"
	"testing"

	"github.com/apache/thrift/lib/go/thrift"

	"github.com/anthony-dong/golang/pkg/bufutils"
	"github.com/anthony-dong/golang/pkg/codec"
	"github.com/anthony-dong/golang/pkg/utils"
)

func TestTestDecodeMessage(t *testing.T) {
	testDecodeMessage(FramedCompact, thrift.CALL, NewTestArgsData(), t)                                                   // test simple request
	testDecodeMessage(FramedCompact, thrift.REPLY, NewTestResultData(), t)                                                //  test simple response
	testDecodeMessage(FramedCompact, thrift.EXCEPTION, thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "错误信息"), t) //  test error msg
	testDecodeMessage(FramedCompact, thrift.ONEWAY, NewTestArgsData(), t)                                                 //  test oneway request
}

func testDecodeMessage(proto Protocol, msgType thrift.TMessageType, msg thrift.TStruct, t *testing.T) {
	buffer := bufutils.NewBuffer()
	encoder := NewTProtocolEncoder(buffer, proto)
	if err := writeThriftMessage(encoder, msgType, msg); err != nil {
		t.Fatal(err)
	}
	t.Log(string(codec.NewBase64Codec().Encode(buffer.Bytes())))
	data, err := DecodeMessage(context.Background(), NewTProtocol(buffer, proto))
	if err != nil {
		t.Fatal(err)
	}
	data.Protocol = proto
	t.Log(utils.ToJson(data))
}

func TestDecode(t *testing.T) {
	proto := UnframedCompact
	data := `gAEAAQAAABJxdWVyeU1ldHJpY0FnZ1N0YXQAAAAADAABDwABDAAAAAELAAEAAAAABAACAAAAAAAAAAALAAMAAAAACAAEAAAAAAAPAAIMAAAAAQsAAQAAAAALAAIAAAAAAA8AAwwAAAABCwABAAAAAAsAAgAAAAAADAAECAABAAAAAQsAAgAAAAALAAMAAAAACwAEAAAAAAAMAAUIAAEAAAABDwACDAAAAAELAAMAAAAACAAEAAAAAQ8ABQsAAAABAAAAAAsABgAAAAAADAAGCAABAAAAAQ8AAgwAAAABCwADAAAAAAgABAAAAAEPAAULAAAAAQAAAAALAAYAAAAAAA8ABwwAAAABCwABAAAAAAgAAgAAAAAPAAMLAAAAAQAAAAAADAAICAABAAAAAAgAAgAAAAAADQAJCwsAAAABAAAAAAAAAAAMAAoPAAEMAAAAAQsAAQAAAAAAAgACAA8AAwsAAAABAAAAAAoABAAAAAAAAAAACwAFAAAAAAgABgAAAAAPAAcIAAAAAQAAAAECAAgACwAJAAAAAAIACgALAAsAAAAACAAMAAAAAAgADf////8MAA4LAAEAAAAACwACAAAAAAAMAA8IAAEAAAABDwACDAAAAAELAAMAAAAACAAEAAAAAQ8ABQsAAAABAAAAAAsABgAAAAAADwAQCwAAAAEAAAAADwARCwAAAAEAAAAADQASCwwAAAABAAAAAAgAAQAAAAELAAIAAAAAAAsAEwAAAAACABQADwAVDAAAAAELAAEAAAAACwACAAAAAAsAAwAAAAAPAAQLAAAAAQAAAAACAAUAAAwAFgsAAQAAAAAADQAXCwsAAAABAAAAAAAAAAAADwALDwAAAAEMAAAAAQsAAQAAAAALAAIAAAAAAAsADAAAAAAMAP8LAAEAAAAACwACAAAAAAsAAwAAAAALAAQAAAAADAAFAgABAAsAAgAAAAAADQAGCwsAAAAFAAAAA2VudgAAABVwcGVfZGV2X2xpdXpoZW5sb25nXzEAAAAEdXNlcgAAAAloZW1pbmdqaW4AAAAKZ2Rwci10b2tlbgAAAwNleUpoYkdjaU9pSlNVekkxTmlJc0luUjVjQ0k2SWtwWFZDSjkuZXlKMlpYSnphVzl1SWpvd0xDSmhkWFJvYjNKcGRIa2lPaUpVUTBVaUxDSndjbWx0WVhKNVFYVjBhRlI1Y0dVaU9pSndjMjBpTENKd2MyMGlPaUpsZUhCc2IzSmxjaTVoY0drdVpYaGxZM1YwYjNJaUxDSjFjMlZ5SWpvaVozVnZlR2x1TG5KcFkyc3pNeUlzSW1WNGNHbHlaVlJwYldVaU9qRTNNVGswT1RBd05qa3NJbVY0ZEdWdWMybHZiaUk2ZXlKamJIVnpkR1Z5WDI1aGJXVWlPaUprWldaaGRXeDBJaXdpYVdSaklqb2lRVWRLVTA1S0lpd2liRzluYVdOaGJGOWpiSFZ6ZEdWeUlqb2laR1ZtWVhWc2RDSXNJbkJvZVhOcFkyRnNYMk5zZFhOMFpYSWlPaUpCWjJkeVpXZGhkR2x2YmlJc0luTmxjblpwWTJWZmRIbHdaU0k2SW1Gd2NGOWxibWRwYm1VaUxDSjZiMjVsSWpvaVEyaHBibUV0UVVjaWZYMC5jaG0wdExPWmJBcFNlOXFPTFdEbVJBM0ZmOXdDUlZ4a2NVT2VwazlUNGZCdVZZV2E4eDZ2QlliN2ZGRVBRc0MxLTFwVWRnZFlFamNCTWdxWlFpbDluY215SmU4TEVSZ09yQjJyRl9MQUxIRU5wS0U0VXRobWZTcVNmOWhwd1pHRkY2b3Mzb3lUY0xxb2ttREpCakhpS3VCZS1PSkhYUGxISEJ5VVprdi0tYTFaMWtpRnJfc19vMHhtWkJLTUppdW85Q09vNERxUlAxWElyMEdRYkNHUWZPOGpSb0l5bEN0aDhscDJoZzRIWF9IRU00cWJEcmJUdUw1NmtpVDhEZ3EtVGhjaWFmWlF5c2F6UThaeVRvVjRKREV0bm84UzVGOHJNWmZJRVFDZVFOVFhHMjY3TGszTUJUZDRZQzVKcHdibmFod2syNzRTcV9NSExYWmZPWTRFckEAAAAAAAAAAAAAAAp1c2VyX2V4dHJhAAAAX3siZGVzdF9zZXJ2aWNlIjoiYWQuc3RhdHMuY3J1eF91bmlmaWVkX3NlcnZpY2UiLCJSUENfUEVSU0lTVF9NT0NLX1RBRyI6InBwZV9kZXZfbGl1emhlbmxvbmdfMSJ9AAAA`
	bData, err := codec.NewBase64Codec().Decode([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
	result, err := DecodeStruct(context.Background(), NewTProtocol(bytes.NewReader(bData), proto))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(utils.ToJson(result, true))
}

func TestFieldOrderMap(t *testing.T) {
	orderMap := FieldOrderMap{}
	orderMap.Set(Field{
		FieldId:   1,
		FieldType: 2,
	}, FieldOrderMap{})

	t.Log(utils.ToJson(orderMap))
}

func TestIsUtf8(t *testing.T) {
	if isValidUTF8(string([]byte{'\u0000'})) {
		t.Log("success")
	}
}
