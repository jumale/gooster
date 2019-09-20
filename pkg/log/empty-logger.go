package log

type EmptyLogger struct{}

func (e EmptyLogger) Debug(v ...interface{}) {}

func (e EmptyLogger) DebugF(msg string, args ...interface{}) {}

func (e EmptyLogger) Info(v ...interface{}) {}

func (e EmptyLogger) InfoF(msg string, args ...interface{}) {}

func (e EmptyLogger) Warn(v ...interface{}) {}

func (e EmptyLogger) WarnF(msg string, args ...interface{}) {}

func (e EmptyLogger) Error(v ...interface{}) {}

func (e EmptyLogger) ErrorF(msg string, args ...interface{}) {}

func (e EmptyLogger) Fatal(v ...interface{}) {}

func (e EmptyLogger) FatalF(msg string, args ...interface{}) {}
