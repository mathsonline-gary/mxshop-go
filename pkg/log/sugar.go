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

func (s *Sugar) Info(args ...interface{}) {
	if s.logger.Level() > LevelInfo {
		return
	}
	_ = s.logger.Log(LevelInfo, s.msgKey, fmt.Sprint(args...))
}

func (s *Sugar) Infof(format string, args ...interface{}) {
	if s.logger.Level() > LevelInfo {
		return
	}
	_ = s.logger.Log(LevelInfo, s.msgKey, fmt.Sprintf(format, args...))
}

func (s *Sugar) Infow(keyvals ...interface{}) {
	if s.logger.Level() > LevelInfo {
		return
	}
	_ = s.logger.Log(LevelInfo, keyvals...)
}

func (s *Sugar) Warn(args ...interface{}) {
	if s.logger.Level() > LevelWarn {
		return
	}
	_ = s.logger.Log(LevelWarn, s.msgKey, fmt.Sprint(args...))
}

func (s *Sugar) Warnf(format string, args ...interface{}) {
	if s.logger.Level() > LevelWarn {
		return
	}
	_ = s.logger.Log(LevelWarn, s.msgKey, fmt.Sprintf(format, args...))
}

func (s *Sugar) Warnw(keyvals ...interface{}) {
	if s.logger.Level() > LevelWarn {
		return
	}
	_ = s.logger.Log(LevelWarn, keyvals...)
}

func (s *Sugar) Error(args ...interface{}) {
	if s.logger.Level() > LevelError {
		return
	}
	_ = s.logger.Log(LevelError, s.msgKey, fmt.Sprint(args...))
}

func (s *Sugar) Errorf(format string, args ...interface{}) {
	if s.logger.Level() > LevelError {
		return
	}
	_ = s.logger.Log(LevelError, s.msgKey, fmt.Sprintf(format, args...))
}

func (s *Sugar) Errorw(keyvals ...interface{}) {
	if s.logger.Level() > LevelError {
		return
	}
	_ = s.logger.Log(LevelError, keyvals...)
}

func (s *Sugar) Fatal(args ...interface{}) {
	if s.logger.Level() > LevelFatal {
		return
	}
	_ = s.logger.Log(LevelFatal, s.msgKey, fmt.Sprint(args...))
}

func (s *Sugar) Fatalf(format string, args ...interface{}) {
	if s.logger.Level() > LevelFatal {
		return
	}
	_ = s.logger.Log(LevelFatal, s.msgKey, fmt.Sprintf(format, args...))
}

func (s *Sugar) Fatalw(keyvals ...interface{}) {
	if s.logger.Level() > LevelFatal {
		return
	}
	_ = s.logger.Log(LevelFatal, keyvals...)
}
