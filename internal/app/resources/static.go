package resources

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/dropwhile/icbt/internal/logger"
)

var (
	//go:embed static
	staticEmbedFs embed.FS
	staticFs      fs.FS
)

type Location int

const (
	Embed Location = iota + 1
	Filesystem
)

func NewStaticFS(loc Location) fs.FS {
	switch loc {
	case Embed:
		var err error
		staticFs, err = fs.Sub(staticEmbedFs, "static")
		if err != nil {
			panic(err)
		}
	case Filesystem:
		sdir := "./internal/app/resources/static/"
		staticFs = os.DirFS(sdir)
	default:
		panic("staticDir must be one of: embed, fs")
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
			slog.DebugContext(r.Context(),
				"cant open file for reading",
				slog.String("filepath", filePath),
				logger.Err(err),
			)
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			slog.DebugContext(r.Context(),
				"cant read file",
				slog.String("filepath", filePath),
				logger.Err(err),
			)
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}
		w.Header().Set("Vary", "Accept-Encoding")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		http.ServeContent(w, r, filePath, time.Time{}, bytes.NewReader(b))
	}
	return fn
}
