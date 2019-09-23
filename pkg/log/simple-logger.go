package log

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rivo/tview"
	"io"
	"runtime"
	"strings"
	"time"
)

func NewSimpleLogger(level Level, target io.Writer) *SimpleLogger {
	return &SimpleLogger{level: level, target: target}
}

type SimpleLogger struct {
	level  Level
	target io.Writer
}

func (l *SimpleLogger) log(level Level, msg string) {
	if level < l.level {
		return
	}
	now := time.Now().Format("15:04:05")
	color := level.Color()
	_, file, line, _ := runtime.Caller(2)
	parts := strings.Split(file, "/")
	prefix := tview.Escape(fmt.Sprintf("%s [%s] %s:%d", now, level, parts[len(parts)-1], line))

	_, err := l.target.Write([]byte(fmt.Sprintf("[%s]%s[-] %s\n", color, prefix, msg)))
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
