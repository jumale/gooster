package events

import (
	"fmt"
	"runtime"
	"strings"
)

func truncateString(str string, num int) string {
	result := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		result = str[0:num] + "..."
	}
	return result
}

func getSubscriberIdFromCaller(skip int) subscriberId {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown"
	}

	parts := strings.Split(file, "/")
	return subscriberId(fmt.Sprintf("%s:%d", parts[len(parts)-1], line))
}

func toString(data interface{}) string {
	switch d := data.(type) {
	case string:
		return d
	case []byte:
		return string(d)
	case rune:
		return string(d)
	default:
		return fmt.Sprintf("%v", d)
	}
}
