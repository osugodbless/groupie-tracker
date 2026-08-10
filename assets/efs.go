package assets

import (
	"embed"
	"io/fs"
	"log"
)

//go:embed "templates" "static"
var files embed.FS

var (
	HTMLFS   = sub(files, "templates")
	StaticFS = sub(files, "static")
)

func sub(f embed.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		log.Fatalf("Failed to create sub FS: %v", err)
	}

	return sub
}
