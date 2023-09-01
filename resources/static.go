package resources

import (
	"embed"
	"io/fs"
	"os"
)

var (
	//go:embed static
	staticEmbedFs embed.FS
	staticFs      fs.FS
)

func NewStaticFS(staticDir string) fs.FS {
	if staticDir == "embed" {
		var err error
		staticFs, err = fs.Sub(staticEmbedFs, "static")
		if err != nil {
			panic(err)
		}
	} else {
		staticFs = os.DirFS(staticDir)
	}

	return staticFs
}
