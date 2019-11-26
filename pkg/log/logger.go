package log

import (
	"encoding/json"
	"strings"
)

type Level int

const (
	Debug Level = 0
	Info        = 1
	Warn        = 2
	Error       = 3
	Fatal       = 4
)

type Logger interface {
	Debug(v ...interface{})
	DebugF(msg string, args ...interface{})

	Info(v ...interface{})
	InfoF(msg string, args ...interface{})

	Warn(v ...interface{})
	WarnF(msg string, args ...interface{})

	Error(v ...interface{})
	ErrorF(msg string, args ...interface{})

	Fatal(v ...interface{})
	FatalF(msg string, args ...interface{})

	Check(err error, msg ...string)
}

func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "FATAL"
	}
}

func (l Level) Color() string {
	switch l {
	case Debug:
		return "gray"
	case Info:
		return "green"
	case Warn:
		return "orange"
	case Error:
		return "red"
	case Fatal:
		return "red"
	default:
		return "red"
	}
}

func (l *Level) UnmarshalJSON(b []byte) error {
	var name string
	if err := json.Unmarshal(b, &name); err != nil {
		return err
	}

	*l = LevelFromString(name)
	return nil
}

func LevelFromString(val string) Level {
	switch strings.ToUpper(val) {
	case "DEBUG":
		return Debug
	case "INFO":
		return Info
	case "WARN":
		return Warn
	case "ERROR":
		return Error
	case "FATAL":
		return Fatal
	default:
		return Warn
	}
}
