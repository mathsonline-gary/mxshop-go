package log

import "strings"

type Level int8

const LevelKey = "level"

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

func (l Level) Key() string {
	return LevelKey
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelPanic:
		return "PANIC"
	case LevelFatal:
		return "FATAL"
	default:
		return ""
	}
}

func ParseLevel(s string) Level {
	switch strings.ToUpper(s) {
	case "DEBUG":
		return LevelDebug
	case "INFO":
		return LevelInfo
	case "WARN":
		return LevelWarn
	case "ERROR":
		return LevelError
	case "PANIC":
		return LevelPanic
	case "FATAL":
		return LevelFatal
	}
	return LevelInfo
}
