package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func JSONRequest(t testing.TB, e *echo.Echo, method, path string, body any, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader
	switch body := body.(type) {
	case nil:
	case string:
		reader = strings.NewReader(body)
	default:
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("encode request body: %v", err)
		}
		reader = bytes.NewReader(data)
	}

	request := httptest.NewRequest(method, path, reader)
	if body != nil {
		request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	for _, cookie := range cookies {
		request.AddCookie(cookie)
	}

	response := httptest.NewRecorder()
	e.ServeHTTP(response, request)
	return response
}

func DecodeJSON(t testing.TB, response *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}

func ResponseCookie(t testing.TB, response *httptest.ResponseRecorder, name string) *http.Cookie {
	t.Helper()
	for _, cookie := range response.Result().Cookies() {
		if cookie.Name == name {
			return cookie
		}
	}
	t.Fatalf("response cookie %q is missing", name)
	return nil
}
