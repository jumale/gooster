package gooster

import "fmt"

func toBytes(data interface{}) []byte {
	switch d := data.(type) {
	case []byte:
		return d
	case string:
		return []byte(d)
	case rune:
		return []byte(string(d))
	default:
		return []byte(fmt.Sprintf("%v", d))
	}
}
