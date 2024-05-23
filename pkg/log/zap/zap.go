package zap

import (
	"fmt"

	"github.com/zycgary/mxshop-go/pkg/log"
	"go.uber.org/zap"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{logger: logger}

}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	var kvlen = len(keyvals)
	if kvlen == 0 || kvlen%2 != 0 {
		l.logger.Warn(fmt.Sprint("key-value must appear in pairs: ", keyvals))
		return nil
	}

	data := make([]zap.Field, 0, (kvlen/2)+1)
	for i := 0; i < kvlen; i += 2 {
		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.logger.Debug("", data...)
	case log.LevelInfo:
		l.logger.Info("", data...)
	case log.LevelWarn:
		l.logger.Warn("", data...)
	case log.LevelError:
		l.logger.Error("", data...)
	case log.LevelPanic:
		l.logger.Panic("", data...)
	case log.LevelFatal:
		l.logger.Fatal("", data...)
	}
	return nil
}

func (l *Logger) Level() log.Level {
	switch l.logger.Level() {
	case zap.DebugLevel:
		return log.LevelDebug
	case zap.InfoLevel:
		return log.LevelInfo
	case zap.WarnLevel:
		return log.LevelWarn
	case zap.ErrorLevel:
		return log.LevelError
	case zap.DPanicLevel, zap.PanicLevel:
		return log.LevelPanic
	case zap.FatalLevel:
		return log.LevelFatal
	default:
		return log.LevelInfo
	}
}

func (l *Logger) Close() error {
	return l.logger.Sync()
}
