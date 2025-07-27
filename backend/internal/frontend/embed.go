package frontend

import (
	"embed"
)

//go:embed dist/browser/**
var FrontendFS embed.FS
