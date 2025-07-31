package logging

import (
	"log/slog"
	"os"
)

type Format int

const (
	FormatText Format = iota
	FormatJSON
	FormatFancy
)

func (f Format) IsValid() bool {
	switch f {
	case FormatText, FormatJSON, FormatFancy:
		return true
	}
	return false
}

func (f Format) NewSLogHandler(options *slog.HandlerOptions) slog.Handler {
	switch f {
	case FormatJSON:
		return slog.NewJSONHandler(os.Stdout, options)
	case FormatFancy:
		return NewFancyHandler(options)
	case FormatText:
		fallthrough
	default:
		return slog.NewTextHandler(os.Stdout, options)
	}
}

func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatFancy:
		return "fancy"
	case FormatText:
		fallthrough
	default:
		return "text"
	}
}

func FormatFromString(s string) Format {
	switch s {
	case "json":
		return FormatJSON
	case "fancy":
		return FormatFancy
	case "text":
		fallthrough
	default:
		return FormatText
	}
}
