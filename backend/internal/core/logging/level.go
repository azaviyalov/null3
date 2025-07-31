package logging

import "log/slog"

type Level int

const (
	LevelDebug Level = iota - 1
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) IsValid() bool {
	switch l {
	case LevelDebug, LevelInfo, LevelWarn, LevelError:
		return true
	}
	return false
}

func (l Level) Level() slog.Level {
	switch l {
	case LevelDebug:
		return slog.LevelDebug
	case LevelWarn:
		return slog.LevelWarn
	case LevelError:
		return slog.LevelError
	case LevelInfo:
		fallthrough
	default:
		return slog.LevelInfo
	}
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelInfo:
		fallthrough
	default:
		return "info"
	}
}

func LevelFromString(s string) Level {
	switch s {
	case "debug":
		return LevelDebug
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	case "info":
		fallthrough
	default:
		return LevelInfo
	}
}
