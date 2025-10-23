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
	fileLine := findExternalSource()
	if fileLine == "" {
		return a
	}
	return slog.Attr{Key: "source", Value: slog.StringValue(fileLine)}
}

func findExternalSource() string {
	pcs := make([]uintptr, 64)
	// skip runtime.Callers + this function + ReplaceAttr etc; start at 3
	n := runtime.Callers(3, pcs)
	if n == 0 {
		return ""
	}
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if frame.File == "" {
			if !more {
				break
			}
			continue
		}
		if isLogFrame(frame) {
			if !more {
				break
			}
			continue
		}
		return trimRepoPrefix(frame.File) + ":" + strconv.Itoa(frame.Line)
	}
	return ""
}

func findExternalFuncName() string {
	pcs := make([]uintptr, 64)
	n := runtime.Callers(3, pcs)
	if n == 0 {
		return ""
	}
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		if frame.Function == "" || frame.File == "" {
			if !more {
				break
			}
			continue
		}
		if isLogFrame(frame) {
			if !more {
				break
			}
			continue
		}
		return trimFunctionName(frame.Function)
	}
	return ""
}

func trimFunctionName(name string) string {
	parts := strings.Split(name, "/")
	return parts[len(parts)-1]
}

func trimRepoPrefix(file string) string {
	marker := "null3/"
	if idx := strings.Index(file, marker); idx != -1 {
		return file[idx+len(marker):]
	}
	return file
}

func isLogFrame(frame runtime.Frame) bool {
	if frame.File != "" {
		if strings.Contains(frame.File, "/log/slog") || strings.Contains(frame.File, "log/slog") {
			return true
		}
		if strings.Contains(frame.File, "backend/internal/core/logging") {
			middlewarePath := "backend/internal/core/logging/middleware.go"
			return frame.File != middlewarePath
		}
	}
	return false
}
