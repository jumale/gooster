package gooster

import (
	"fmt"
	"github.com/rivo/tview"
	"time"
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

// -------------------------------------------------- //

type LogLevel int

const (
	LogDebug LogLevel = 0
	LogInfo           = 1
	LogWarn           = 2
	LogError          = 3
	LogFatal          = 4
)

func logLevelName(l LogLevel) string {
	switch l {
	case LogDebug:
		return "DEBUG"
	case LogInfo:
		return "INFO"
	case LogWarn:
		return "WARN"
	case LogError:
		return "ERROR"
	default:
		return "FATAL"
	}
}

func logLevelColor(l LogLevel) string {
	switch l {
	case LogDebug:
		return "gray"
	case LogInfo:
		return "green"
	case LogWarn:
		return "orange"
	case LogError:
		return "red"
	default:
		return "red"
	}
}

// -------------------------------------------------- //

func NewSelfLogger(level LogLevel, em *EventManager) *SelfLogger {
	return &SelfLogger{em: em, level: level}
}

type SelfLogger struct {
	em    *EventManager
	level LogLevel
}

func (l *SelfLogger) send(level LogLevel, msg string) {
	if level < l.level {
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	color := logLevelColor(level)
	prefix := tview.Escape(fmt.Sprintf("[%s] [%s]", now, logLevelName(level)))

	l.em.Dispatch(Event{
		Id:   EventOutputMessage,
		Data: fmt.Sprintf("[%s]%s[-] %s", color, prefix, msg),
	})
}

func (l *SelfLogger) Debug(v ...interface{}) {
	l.send(LogDebug, fmt.Sprint(v...))
}

func (l *SelfLogger) DebugF(msg string, args ...interface{}) {
	l.send(LogDebug, fmt.Sprintf(msg, args...))
}

func (l *SelfLogger) Info(v ...interface{}) {
	l.send(LogInfo, fmt.Sprint(v...))
}

func (l *SelfLogger) InfoF(msg string, args ...interface{}) {
	l.send(LogInfo, fmt.Sprintf(msg, args...))
}

func (l *SelfLogger) Warn(v ...interface{}) {
	l.send(LogWarn, fmt.Sprint(v...))
}

func (l *SelfLogger) WarnF(msg string, args ...interface{}) {
	l.send(LogWarn, fmt.Sprintf(msg, args...))
}

func (l *SelfLogger) Error(v ...interface{}) {
	l.send(LogError, fmt.Sprint(v...))
}

func (l *SelfLogger) ErrorF(msg string, args ...interface{}) {
	l.send(LogError, fmt.Sprintf(msg, args...))
}

func (l *SelfLogger) Fatal(v ...interface{}) {
	l.send(LogFatal, fmt.Sprint(v...))
}

func (l *SelfLogger) FatalF(msg string, args ...interface{}) {
	l.send(LogFatal, fmt.Sprintf(msg, args...))
}
