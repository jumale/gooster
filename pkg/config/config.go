package config

type Reader interface {
	Read(jsonPath string, target interface{}) error
}
