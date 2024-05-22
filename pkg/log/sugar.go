package log

import (
	"fmt"
)

type Sugar struct {
	logger Logger
	msgKey string
}

func NewSugar(logger Logger) *Sugar {
	return &Sugar{logger: logger}
}

func (s *Sugar) Log(level Level, keyvals ...interface{}) {
	_ = s.logger.Log(level, keyvals...)
}

func (s *Sugar) Debug(args ...interface{}) {
	if s.logger.Level() > LevelDebug {
		return
	}
	_ = s.logger.Log(LevelDebug, s.msgKey, fmt.Sprint(args...))
}

func (s *Sugar) Debugf(format string, args ...interface{}) {
	if s.logger.Level() > LevelDebug {
		return
	}
	_ = s.logger.Log(LevelDebug, s.msgKey, fmt.Sprintf(format, args...))
}

func (s *Sugar) Debugw(keyvals ...interface{}) {
	if s.logger.Level() > LevelDebug {
		return
	}
	_ = s.logger.Log(LevelDebug, keyvals...)
}
