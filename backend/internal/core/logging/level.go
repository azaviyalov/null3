package logging

import "log/slog"

type Level uint8

const (
	LevelDefault Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
)

func (l Level) IsValid() bool {
	switch l {
	case LevelDefault, LevelDebug, LevelInfo, LevelWarn, LevelError:
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
	case LevelDefault, LevelInfo:
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
	case LevelDefault, LevelInfo:
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
