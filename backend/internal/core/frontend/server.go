package frontend

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRoutes(e *echo.Echo, config Config) {
	if !config.EnableFrontendDist {
		e.Logger.Info("Frontend dist is disabled, skipping static file serving")
		return
	}

	frontendFS, err := fs.Sub(FrontendFS, "fs")
	if err != nil {
		panic("failed to sub fs: " + err.Error())
	}

	// Traverse and patch files
	patchedFiles := make(map[string][]byte)
	fs.WalkDir(frontendFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		f, err := frontendFS.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()
		content, err := io.ReadAll(f)
		if err != nil {
			return nil
		}

		// Patch API URL placeholder
		if bytes.Contains(content, []byte("%%API_URL%%")) {
			content = bytes.ReplaceAll(content, []byte("%%API_URL%%"), []byte(config.APIURL))
			patchedFiles[path] = content
		}
		return nil
	})

	memfs := &memFS{
		files: patchedFiles,
		orig:  frontendFS,
	}

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		HTML5:      true,
		Filesystem: http.FS(memfs),
	}))
}
