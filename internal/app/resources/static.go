package resources

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"
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

func FileServer(fsys fs.FS, stripPrefix string) http.HandlerFunc {
	staticServer := http.FileServer(http.FS(fsys))
	if stripPrefix != "" {
		staticServer = http.StripPrefix(stripPrefix, staticServer)
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		staticServer.ServeHTTP(w, r)
	}
	return fn
}

func ServeSingle(fsys fs.FS, filePath string) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		f, err := fsys.Open(filePath)
		if err != nil {
			log.Debug().Err(err).
				Str("filepath", filePath).
				Msg("cant open file for reading")
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			log.Debug().Err(err).
				Str("filepath", filePath).
				Msg("cant read file")
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeContent(w, r, filePath, time.Time{}, bytes.NewReader(b))
	}
	return fn
}