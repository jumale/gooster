package gooster

import (
	"fmt"
	"os"
)

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

func getWd() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}
