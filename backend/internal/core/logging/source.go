package logging

import (
	"log/slog"
	"runtime"
	"strconv"
	"strings"
)

func ReplaceSourceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key != "source" {
		return a
	}
	source, ok := a.Value.Any().(*slog.Source)
	if !ok || source == nil {
		return a
	}
	return slog.String("source", formatSource(source.File, source.Line))
}

func sourceFromPC(pc uintptr) string {
	if pc == 0 {
		return ""
	}
	frame, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	if frame.File == "" {
		return ""
	}
	return formatSource(frame.File, frame.Line)
}

func formatSource(file string, line int) string {
	return trimRepoPrefix(file) + ":" + strconv.Itoa(line)
}

func trimRepoPrefix(file string) string {
	marker := "null3/"
	if idx := strings.Index(file, marker); idx != -1 {
		return file[idx+len(marker):]
	}
	return file
}
