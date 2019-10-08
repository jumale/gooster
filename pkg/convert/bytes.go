package convert

import "fmt"

func ToBytes(data interface{}) []byte {
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

func ToString(data interface{}) string {
	switch d := data.(type) {
	case []byte:
		return string(d)
	case rune:
		return string(d)
	case string:
		return d
	default:
		return fmt.Sprintf("%v", d)
	}
}
