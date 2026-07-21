package logging

import (
	"io"
	"log/slog"
)

type Format int

const (
	FormatText Format = iota
	FormatJSON
	FormatFancy
)

func (f Format) NewSLogHandler(writer io.Writer, options *slog.HandlerOptions) slog.Handler {
	switch f {
	case FormatJSON:
		return slog.NewJSONHandler(writer, options)
	case FormatFancy:
		return NewFancyHandler(writer, options)
	case FormatText:
		fallthrough
	default:
		return slog.NewTextHandler(writer, options)
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
