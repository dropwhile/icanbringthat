package xhandler

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/internal/util"
	"github.com/dropwhile/icbt/resources"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type XHandler struct {
	Db      model.PgxHandle
	Tpl     resources.TemplateMap
	SessMgr *session.SessionMgr
	Mailer  util.MailSender
	Hmac    *util.Hmac
}

func (x *XHandler) Template(name string) (*template.Template, error) {
	return x.Tpl.Get(name)
}

func (x *XHandler) TemplateExecute(w io.Writer, name string, vars map[string]any) error {
	tpl, err := x.Tpl.Get(name)
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

func (x *XHandler) TemplateExecuteSub(w io.Writer, name, subname string, vars map[string]any) error {
	tpl, err := x.Tpl.Get(name)
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

func (x *XHandler) NotFound(w http.ResponseWriter, r *http.Request) {
	x.Error(w, "Not Found", 404)
}

func (x *XHandler) Error(w http.ResponseWriter, statusMsg string, code int) {
	w.Header().Set("content-type", "text/html")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	err := x.TemplateExecute(w, "error-page.gohtml", map[string]any{
		"ErrorCode":   404,
		"ErrorStatus": statusMsg,
		"title":       fmt.Sprintf("%d - %s", code, statusMsg),
	})
	if err != nil {
		// error rendering template, so just return a very basic status page
		log.Debug().Err(err).Msg("custom error status page render issue")
		fmt.Fprintln(w, statusMsg)
		return
	}
}