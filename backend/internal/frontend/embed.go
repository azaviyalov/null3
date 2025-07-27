package frontend

import (
	"embed"
)

//go:embed fs/**
var FrontendFS embed.FS
