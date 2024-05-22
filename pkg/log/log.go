package log

// Logger is the interface that wraps the basic Log method.
type Logger interface {
	Log(level Level, keyvals ...interface{}) error
	Level() Level
}

func New(opt *Option) Logger {
	return nil
}
