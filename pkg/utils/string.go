package utils

import (
	"bytes"
	"strconv"
	"strings"
)

// EscapeQuote escape the string the single quote, double quote, and backtick
func EscapeQuote(str string) string {
	type Escape struct {
		From string
		To   string
	}
	escape := []Escape{
		{From: "`", To: ""}, // remove the backtick
		{From: `\`, To: `\\`},
		{From: `'`, To: `\'`},
		{From: `"`, To: `\"`},
	}

	for _, e := range escape {
		str = strings.ReplaceAll(str, e.From, e.To)
	}
	return str
}

// IsEmpty 是否是空字符串
func IsEmpty(s string) bool {
	if s == "" {
		return true
	}

	return strings.TrimSpace(s) == ""
}

// ConcatString 连接字符串
// NOTE: 性能比fmt.Sprintf和+号要好
func ConcatString(s ...string) string {
	if len(s) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for _, i := range s {
		buffer.WriteString(i)
	}
	return buffer.String()
}

// ConcatStringBySlash concat string by slash
func ConcatStringBySlash(s ...string) string {
	var buffer bytes.Buffer
	for idx, i := range s {
		buffer.WriteString(i)
		if idx != len(s)-1 {
			buffer.WriteString("/")
		}
	}
	return buffer.String()
}

// StringToUint64 字符串转uint64
func StringToUint64(str string) (uint64, error) {
	if str == "" {
		return 0, nil
	}
	valInt, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return uint64(valInt), nil
}

// StringToInt64 字符串转int64
func StringToInt64(str string) (int64, error) {
	if str == "" {
		return 0, nil
	}
	valInt, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return int64(valInt), nil
}

// StringToInt 字符串转int
func StringToInt(str string) (int, error) {
	if str == "" {
		return 0, nil
	}
	valInt, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}

	return valInt, nil
}
