package utils

import (
	"strconv"
	"time"
)

const (
	FormatTimeV1 = "2006-01-02 15:04:05"
	FormatTimeV2 = "2006/1-2"
	FormatTimeV3 = "2006-01-02 15:04:05.000"
)

// TimeToSeconds 时间之差 s 输出 0.100010s.
func TimeToSeconds(duration time.Duration) string {
	// 1s=1000ms 1ms=1000us  保留6位到us
	return strconv.FormatInt(int64(duration/time.Second), 10) + "s"
}

// Float642String 除固定值，保留固定小数位.
func Float642String(num float64, saveDecimalPoint int) string {
	return strconv.FormatFloat(num, 'f', saveDecimalPoint, 64)
}

type JsonDuration time.Duration

func (j JsonDuration) Duration() time.Duration {
	return time.Duration(j)
}

func (j JsonDuration) String() string {
	return j.Duration().String()
}

func (j JsonDuration) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(time.Duration(j).String())), nil
}

func NewJsonDuration(duration time.Duration) JsonDuration {
	return JsonDuration(duration)
}

func (j *JsonDuration) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		unquote, err := strconv.Unquote(string(data))
		if err != nil {
			return err
		}
		if unquote == "" {
			*j = 0
			return nil
		}
		duration, err := time.ParseDuration(unquote)
		if err != nil {
			return err
		}
		*j = JsonDuration(duration)
		return nil
	}
	if len(data) == 4 && string(data) == "null" {
		*j = 0
		return nil
	}
	duration, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*j = JsonDuration(duration)
	return nil
}
