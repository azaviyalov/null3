package frontend

import (
	"bytes"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitModule(e *echo.Echo) {
	frontendFS, err := fs.Sub(FrontendFS, "dist/browser")
	if err != nil {
		panic("failed to sub dist/browser: " + err.Error())
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
		apiURL := os.Getenv("API_URL")
		if bytes.Contains(content, []byte("%%API_URL%%")) {
			content = bytes.ReplaceAll(content, []byte("%%API_URL%%"), []byte(apiURL))
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
