package log

type EmptyLogger struct{}

func (e EmptyLogger) Debug(_ ...interface{}) {}

func (e EmptyLogger) DebugF(_ string, _ ...interface{}) {}

func (e EmptyLogger) Info(_ ...interface{}) {}

func (e EmptyLogger) InfoF(_ string, _ ...interface{}) {}

func (e EmptyLogger) Warn(v ...interface{}) {}

func (e EmptyLogger) WarnF(_ string, _ ...interface{}) {}

func (e EmptyLogger) Error(_ ...interface{}) {}

func (e EmptyLogger) ErrorF(_ string, _ ...interface{}) {}

func (e EmptyLogger) Fatal(_ ...interface{}) {}

func (e EmptyLogger) FatalF(_ string, _ ...interface{}) {}

func (e EmptyLogger) Check(_ error, _ ...string) {}
