package frontend

import (
	"bytes"
	"io/fs"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type memFS struct {
	files map[string][]byte
	orig  fs.FS
}

func (m *memFS) Open(name string) (fs.File, error) {
	cleanName := path.Clean(strings.TrimPrefix(name, "/"))
	if data, ok := m.files[cleanName]; ok {
		return &memFile{
			reader: bytes.NewReader(data),
			name:   cleanName,
			size:   int64(len(data)),
		}, nil
	}
	return m.orig.Open(cleanName)
}

type memFile struct {
	reader *bytes.Reader
	name   string
	size   int64
}

func (f *memFile) Stat() (fs.FileInfo, error) {
	return &memFileInfo{name: f.name, size: f.size}, nil
}
func (f *memFile) Read(p []byte) (int, error) { return f.reader.Read(p) }
func (f *memFile) Close() error               { return nil }
func (f *memFile) Seek(offset int64, whence int) (int64, error) {
	return f.reader.Seek(offset, whence)
}

type memFileInfo struct {
	name string
	size int64
}

func (fi *memFileInfo) Name() string       { return filepath.Base(fi.name) }
func (fi *memFileInfo) Size() int64        { return fi.size }
func (fi *memFileInfo) Mode() fs.FileMode  { return 0444 }
func (fi *memFileInfo) ModTime() time.Time { return time.Time{} }
func (fi *memFileInfo) IsDir() bool        { return false }
func (fi *memFileInfo) Sys() any           { return nil }
