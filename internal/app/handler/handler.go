package handler

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/cactus/mlog"
	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/resources"
)

type Handler struct {
	Db      *model.DB
	Tpl     resources.TemplateMap
	SessMgr *session.SessionMgr
}

func (h *Handler) Template(name string) (*template.Template, error) {
	return h.Tpl.Get(name)
}

func (h *Handler) TemplateExecute(w io.Writer, name string, vars map[string]any) error {
	tpl, err := h.Tpl.Get(name)
	if err != nil {
		mlog.Infom("template locate error", mlog.Map{"tpl": name, "err": err, "vars": vars})
		return err
	}
	err = tpl.Execute(w, vars)
	if err != nil {
		mlog.Infom("template locate error", mlog.Map{"tpl": name, "err": err})
		return err
	}
	return nil
}

func (h *Handler) TemplateExecuteSub(w io.Writer, name, subname string, vars map[string]any) error {
	tpl, err := h.Tpl.Get(name)
	if err != nil {
		mlog.Infom("template locate error", mlog.Map{"tpl": name, "sub": subname, "err": err, "vars": vars})
		return err
	}
	err = tpl.ExecuteTemplate(w, subname, vars)
	if err != nil {
		mlog.Infom("template locate error", mlog.Map{"tpl": name, "sub": subname, "err": err})
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
	r.ParseForm()
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
	err := h.TemplateExecute(w, x, tplVars)
	if err != nil {
		fmt.Fprint(w, err)
	}
}
