package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

//err非空直接进入panic模式。
func CheckError(err error, message string) {
	if err != nil {
		log.Panicln(message, err)
	}
}

//err非空直接进入panic模式。
func CheckErrorf(err error, format string, v ...interface{}) {
	if err != nil {
		message := fmt.Sprintf(format, v)
		CheckError(err, message)
	}
}

//err非空则显示错误并继续执行。
func ShowError(err error, message string) {
	if err != nil {
		log.Println(message, err)
	}
}

//err非空则显示错误并继续执行。
func ShowErrorf(err error, format string, v ...interface{}) {
	if err != nil {
		message := fmt.Sprintf(format, v)
		CheckError(err, message)
	}
}

//将v作为JSON字符串，解析到指针d中，若v为空，则将d指向的结构生成为空JSON。
func ToJson(v interface{}, d interface{}) error {
	if v != nil {
		s := v.(string)
		if strings.TrimSpace(s) == "" {
			return json.Unmarshal([]byte("{}"), d)
		} else {
			return json.Unmarshal([]byte(s), d)
		}
	}
	return json.Unmarshal([]byte("{}"), d)
}

//从v解析int值，若v为空则返回0。
func ParseInt(v interface{}) int {
	if v == nil {
		return 0
	}
	i := 0
	switch v.(type) {
	case string:
		{
			i, _ = strconv.Atoi(v.(string))
		}
	case float64:
		{
			i = int(v.(float64))
		}
	case int:
		{
			i = v.(int)
		}
	}
	return i
}

//从v解析bool值，若v为空则返回false。
func ParseBool(v interface{}) bool {
	if v == nil {
		return false
	}
	switch v.(type) {
	case string:
		{
			return v.(string) == "true"
		}
	}
	return false
}

//从v解析string值，若v为空则返回空字符串。
func ParseString(v interface{}) string {
	if v == nil {
		return ""
	}
	s := v.(string)
	if strings.TrimSpace(s) == "" {
		return ""
	}
	return s
}

//从v解析string值，若v为空白则返回空def。
func ParseStringDef(v interface{}, def string) string {
	if v == nil {
		return def
	}
	s := v.(string)
	if strings.TrimSpace(s) == "" {
		return def
	}
	return s
}

//从v解析string值，若v非空则赋值给dest。
func ParseStringAndSetWhenNotBlank(dest *string, v interface{}) {
	if v != nil {
		*dest = v.(string)
	}
}

//从v解析int值，若v非0则赋值给dest。
func ParseIntAndSetWhenNotZero(dest *int, v interface{}) {
	if v != nil {
		*dest = ParseInt(v)
	}
}

//从v解析string值，若v非0，则比较该字符串，如果等于“on”则将true赋值给b，否则将false赋值给b。
func ParseChannelStatusAndSet(b *bool, v interface{}) {
	s := ParseString(v)
	if s != "" {
		*b = s == "on"
	}
}

//从v解析string值，若v非0，则比较该字符串，如果不等于“close”则将true赋值给b，否则将false赋值给b。
func ParseDoorMagnetStatusAndSet(b *bool, v interface{}) {
	s := ParseString(v)
	if s != "" {
		*b = s != "close"
	}
}
