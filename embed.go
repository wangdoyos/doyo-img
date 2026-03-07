package main

import (
	"embed"
	"io/fs"
)

//go:embed all:web/dist
var frontendDist embed.FS

func FrontendFS() fs.FS {
	sub, err := fs.Sub(frontendDist, "web/dist")
	if err != nil {
		return nil
	}
	return sub
}
