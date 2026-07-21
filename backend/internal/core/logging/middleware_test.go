package logging_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/azaviyalov/null3/backend/internal/core/logging"
	"github.com/labstack/echo/v4"
)

func TestRequestLoggerLogsSuccessfulRequest(t *testing.T) {
	logBuffer := installJSONLogger(t)
	e := echo.New()
	e.Use(logging.RequestLogger())
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("X-Request-Id", "request-123")
	response := httptest.NewRecorder()

	e.ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusNoContent)
	}
	if got := response.Header().Get("X-Request-Id"); got != "request-123" {
		t.Fatalf("response request ID = %q, want %q", got, "request-123")
	}
	record := decodeLogRecord(t, logBuffer)
	assertLogField(t, record, "level", "INFO")
	assertLogField(t, record, "msg", "HTTP request completed")
	assertLogField(t, record, "request_id", "request-123")
	assertLogField(t, record, "method", http.MethodGet)
	assertLogField(t, record, "path", "/health")
	assertLogField(t, record, "status", float64(http.StatusNoContent))
	assertDurationField(t, record, "latency")
}

func TestRequestLoggerUsesEchoRequestID(t *testing.T) {
	logBuffer := installJSONLogger(t)
	e := echo.New()
	e.Use(logging.RequestLogger())
	e.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	response := httptest.NewRecorder()

	e.ServeHTTP(response, request)

	requestID := response.Header().Get("X-Request-Id")
	if requestID == "" {
		t.Fatal("response request ID is empty")
	}
	record := decodeLogRecord(t, logBuffer)
	assertLogField(t, record, "request_id", requestID)
}

func TestRequestLoggerUsesRoutePattern(t *testing.T) {
	logBuffer := installJSONLogger(t)
	e := echo.New()
	e.Use(logging.RequestLogger())
	e.GET("/invites/:token", func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	request := httptest.NewRequest(http.MethodGet, "/invites/private-token", nil)
	response := httptest.NewRecorder()

	e.ServeHTTP(response, request)

	record := decodeLogRecord(t, logBuffer)
	assertLogField(t, record, "path", "/invites/:token")
	if bytes.Contains(logBuffer.Bytes(), []byte("private-token")) {
		t.Fatal("request log contains a route parameter value")
	}
}

func TestRequestLoggerLogsStreamingResponseOnce(t *testing.T) {
	loggedBeforeCompletion := false
	logBuffer := installJSONLogger(t)
	e := echo.New()
	e.Use(logging.RequestLogger())
	e.GET("/stream", func(c echo.Context) error {
		if _, err := c.Response().Write([]byte("first")); err != nil {
			return err
		}
		if logBuffer.Len() != 0 {
			loggedBeforeCompletion = true
		}
		_, err := c.Response().Write([]byte("second"))
		if logBuffer.Len() != 0 {
			loggedBeforeCompletion = true
		}
		return err
	})
	request := httptest.NewRequest(http.MethodGet, "/stream", nil)
	response := httptest.NewRecorder()

	e.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}
	if response.Body.String() != "firstsecond" {
		t.Fatalf("body = %q, want %q", response.Body.String(), "firstsecond")
	}
	if loggedBeforeCompletion {
		t.Fatal("request was logged before the handler completed")
	}
	record := decodeLogRecord(t, logBuffer)
	assertLogField(t, record, "msg", "HTTP request completed")
	assertLogField(t, record, "status", float64(http.StatusOK))
}

func TestRequestLoggerLogsHandlerErrors(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		wantLevel   string
		wantMessage string
	}{
		{
			name:        "client error",
			status:      http.StatusBadRequest,
			wantLevel:   "WARN",
			wantMessage: "HTTP request completed with non-fatal error",
		},
		{
			name:        "server error",
			status:      http.StatusInternalServerError,
			wantLevel:   "ERROR",
			wantMessage: "HTTP request completed with error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logBuffer := installJSONLogger(t)
			e := echo.New()
			e.Use(logging.RequestLogger())
			e.GET("/failure", func(echo.Context) error {
				return echo.NewHTTPError(tt.status, "request failed")
			})
			request := httptest.NewRequest(http.MethodGet, "/failure", nil)
			response := httptest.NewRecorder()

			e.ServeHTTP(response, request)

			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d", response.Code, tt.status)
			}
			record := decodeLogRecord(t, logBuffer)
			assertLogField(t, record, "level", tt.wantLevel)
			assertLogField(t, record, "msg", tt.wantMessage)
			assertLogField(t, record, "status", float64(tt.status))
			if _, ok := record["error"]; !ok {
				t.Fatal("log record does not contain error field")
			}
		})
	}
}

func installJSONLogger(t *testing.T) *bytes.Buffer {
	t.Helper()

	previous := slog.Default()
	buffer := &bytes.Buffer{}
	slog.SetDefault(slog.New(slog.NewJSONHandler(buffer, nil)))
	t.Cleanup(func() {
		slog.SetDefault(previous)
	})
	return buffer
}

func decodeLogRecord(t *testing.T, buffer *bytes.Buffer) map[string]any {
	t.Helper()

	var record map[string]any
	decoder := json.NewDecoder(buffer)
	if err := decoder.Decode(&record); err != nil {
		t.Fatalf("decode log record: %v", err)
	}
	var extra json.RawMessage
	if err := decoder.Decode(&extra); err != io.EOF {
		t.Fatalf("log buffer contains unexpected trailing data: %v", err)
	}
	return record
}

func assertLogField(t *testing.T, record map[string]any, key string, want any) {
	t.Helper()

	if got := record[key]; got != want {
		t.Errorf("log field %q = %v, want %v", key, got, want)
	}
}

func assertDurationField(t *testing.T, record map[string]any, key string) {
	t.Helper()

	value, ok := record[key].(string)
	if !ok {
		t.Fatalf("log field %q = %v, want duration string", key, record[key])
	}
	if _, err := time.ParseDuration(value); err != nil {
		t.Fatalf("log field %q = %q, want duration: %v", key, value, err)
	}
}
