package handler

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/resources"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	Db      model.PgxHandle
	Tpl     resources.TemplateMap
	SessMgr *session.SessionMgr
	Mailer  *util.Mailer
	Hmac    *util.Hmac
}

/*
type WrappableHandler func(*Handler, http.ResponseWriter, *http.Request)

func (h *Handler) Wrap(wrapped WrappableHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wrapped(h, w, r)
	}
}
*/

func (h *Handler) Template(name string) (*template.Template, error) {
	return h.Tpl.Get(name)
}

func (h *Handler) HxCurrentUrlHasPrefix(r *http.Request, prefix string) bool {
	htmxCurrentUrl := r.Header.Get("HX-Current-URL")
	if htmxCurrentUrl != "" {
		u, err := url.Parse(htmxCurrentUrl)
		return err == nil && strings.HasPrefix(u.Path, prefix)
	}
	return false
}

func (h *Handler) TemplateExecute(w io.Writer, name string, vars map[string]any) error {
	tpl, err := h.Tpl.Get(name)
	if err != nil {
		log.Info().
			Err(err).
			Str("tpl", name).
			Msg("template locate error")
		return err
	}
	err = tpl.Execute(w, vars)
	if err != nil {
		log.Info().
			Err(err).
			Str("tpl", name).
			Dict("vars", zerolog.Dict().Fields(vars)).
			Msg("template execute error")
		return err
	}
	return nil
}

func (h *Handler) TemplateExecuteSub(w io.Writer, name, subname string, vars map[string]any) error {
	tpl, err := h.Tpl.Get(name)
	if err != nil {
		log.Info().
			Err(err).
			Str("tpl", name).
			Str("sub", subname).
			Msg("template locate error")
		return err
	}
	err = tpl.ExecuteTemplate(w, subname, vars)
	if err != nil {
		log.Info().
			Err(err).
			Str("tpl", name).
			Str("sub", subname).
			Dict("vars", zerolog.Dict().Fields(vars)).
			Msg("template execute error")
		return err
	}
	return nil
}

// SetHeader is a convenience handler to set a response header key/value
func SetHeader(key, value string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(key, value)
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (h *Handler) TestTemplates(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	x := r.Form.Get("tpl")
	y := r.Form.Get("nav")
	if x == "" {
		x = "index.gohtml"
	}
	if y == "" {
		y = "dashboard"
	}
	tplVars := map[string]any{
		"user":  "nope",
		"title": y,
		"nav":   y,
	}
	err = h.TemplateExecute(w, x, tplVars)
	if err != nil {
		fmt.Fprint(w, err)
	}
}
