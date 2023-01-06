package logger

type NullLogger struct{}

func NewNullLogger() *NullLogger {
	return &NullLogger{}
}

func (l *NullLogger) Debug(args ...interface{}) {}

func (l *NullLogger) Debugf(format string, args ...interface{}) {}

func (l *NullLogger) Info(args ...interface{}) {}

func (l *NullLogger) Infof(format string, args ...interface{}) {}

func (l *NullLogger) Warning(args ...interface{}) {}

func (l *NullLogger) Warningf(format string, args ...interface{}) {}

func (l *NullLogger) Error(args ...interface{}) {}

func (l *NullLogger) Errorf(format string, args ...interface{}) {}

func (l *NullLogger) Panic(args ...interface{}) {}

func (l *NullLogger) Panicf(format string, args ...interface{}) {}

func (l *NullLogger) Fatal(args ...interface{}) {}

func (l *NullLogger) Fatalf(format string, args ...interface{}) {}

func (l *NullLogger) SetExtraCallDepth(depth int) {}

func (l *NullLogger) ResetExtraCallDepth() {}
