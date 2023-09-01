package api

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type API struct {
	db  *model.DB
	tpl *template.Template
	*chi.Mux
}

func New(db *model.DB, tpl *template.Template) *API {
	r := chi.NewRouter()
	if mlog.HasDebug() {
		r.Use(middleware.Logger)
	}
	r.Use(middleware.RedirectSlashes)
	r.Use(middleware.GetHead)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	r.Get("/testtemplate", func(w http.ResponseWriter, r *http.Request) {
		ctx := map[string]string{
			"path": r.URL.Path,
			"year": fmt.Sprintf("%d", time.Now().Local().Year()),
		}
		tpl.ExecuteTemplate(w, "index.html", ctx)
	})

	return &API{db: db, tpl: tpl, Mux: r}
}
