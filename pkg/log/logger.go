package log

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
}

func LevelName(l Level) string {
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

func LevelColor(l Level) string {
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
