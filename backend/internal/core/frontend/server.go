package frontend

import (
	"bytes"
	"io/fs"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRoutes(e *echo.Echo, config Config) {
	if !config.EnableFrontendDist {
		slog.Info("frontend dist is disabled, skipping static file serving")
		return
	}

	frontendFS, err := fs.Sub(FrontendFS, "fs")
	if err != nil {
		panic("failed to sub fs: " + err.Error())
	}

	registerStaticRoutes(e, frontendFS, config.APIURL)
}

func registerStaticRoutes(e *echo.Echo, frontendFS fs.FS, apiURL string) {
	patchedFiles := make(map[string][]byte)
	err := fs.WalkDir(frontendFS, ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		content, err := fs.ReadFile(frontendFS, path)
		if err != nil {
			return err
		}

		if bytes.Contains(content, []byte("%%API_URL%%")) {
			content = bytes.ReplaceAll(content, []byte("%%API_URL%%"), []byte(apiURL))
			patchedFiles[path] = content
		}
		return nil
	})
	if err != nil {
		panic("patch frontend fs: " + err.Error())
	}

	frontendFS = &memFS{
		files: patchedFiles,
		orig:  frontendFS,
	}
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Filesystem: http.FS(frontendFS),
	}))
}
