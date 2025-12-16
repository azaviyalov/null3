package main

import (
	"github.com/azaviyalov/null3/backend/internal/app"
)

func main() {
	app := app.New()
	app.Start()
}
