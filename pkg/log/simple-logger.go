package log

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"runtime"
	"strings"
	"time"
)

type SimpleLoggerConfig struct {
	Level  Level
	Format string
}

const defaultLogFormat = "[<color>]<time> [<level>] <caller>[-] <msg>\n"

func NewSimpleLogger(target io.Writer, cfg SimpleLoggerConfig) *SimpleLogger {
	if cfg.Format == "" {
		cfg.Format = defaultLogFormat
	}
	return &SimpleLogger{target: target, cfg: cfg}
}

type SimpleLogger struct {
	cfg    SimpleLoggerConfig
	target io.Writer
}

func (l *SimpleLogger) log(level Level, msg string) {
	if level < l.cfg.Level {
		return
	}
	now := time.Now().Format("15:04:05")
	color := level.Color()

	_, file, line, _ := runtime.Caller(2)
	parts := strings.Split(file, "/")
	caller := fmt.Sprintf("%s:%d", parts[len(parts)-1], line)

	log := l.cfg.Format
	log = strings.Replace(log, "<color>", color, -1)
	log = strings.Replace(log, "<time>", now, -1)
	log = strings.Replace(log, "<level>", level.String(), -1)
	log = strings.Replace(log, "<caller>", caller, -1)
	log = strings.Replace(log, "<msg>", msg, -1)

	_, err := l.target.Write([]byte(log))
	if err != nil {
		panic(errors.WithMessage(err, "writing to log target"))
	}
}

func (l *SimpleLogger) Debug(v ...interface{}) {
	l.log(Debug, fmt.Sprint(v...))
}

func (l *SimpleLogger) DebugF(msg string, args ...interface{}) {
	l.log(Debug, fmt.Sprintf(msg, args...))
}

func (l *SimpleLogger) Info(v ...interface{}) {
	l.log(Info, fmt.Sprint(v...))
}

func (l *SimpleLogger) InfoF(msg string, args ...interface{}) {
	l.log(Info, fmt.Sprintf(msg, args...))
}

func (l *SimpleLogger) Warn(v ...interface{}) {
	l.log(Warn, fmt.Sprint(v...))
}

func (l *SimpleLogger) WarnF(msg string, args ...interface{}) {
	l.log(Warn, fmt.Sprintf(msg, args...))
}

func (l *SimpleLogger) Error(v ...interface{}) {
	l.log(Error, fmt.Sprint(v...))
}

func (l *SimpleLogger) ErrorF(msg string, args ...interface{}) {
	l.log(Error, fmt.Sprintf(msg, args...))
}

func (l *SimpleLogger) Fatal(v ...interface{}) {
	l.log(Fatal, fmt.Sprint(v...))
}

func (l *SimpleLogger) FatalF(msg string, args ...interface{}) {
	l.log(Fatal, fmt.Sprintf(msg, args...))
}

func (l *SimpleLogger) Check(err error, msgAndArgs ...string) {
	if err == nil {
		return
	}
	if len(msgAndArgs) > 0 {
		err = errors.WithMessage(err, fmt.Sprintf(msgAndArgs[0], msgAndArgs[1:]))
	}
	l.Error(err)
}
