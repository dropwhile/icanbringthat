package xhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/crypto"
	"github.com/dropwhile/icbt/internal/mail"
	"github.com/dropwhile/icbt/internal/session"
	"github.com/dropwhile/icbt/resources"
)

type XHandler struct {
	Db          model.PgxHandle
	Redis       *redis.Client
	TemplateMap resources.TemplateMap
	SessMgr     *session.SessionMgr
	Mailer      mail.MailSender
	MAC         *crypto.MAC
	BaseURL     string
	IsProd      bool
}

func (x *XHandler) Template(name string) (resources.TemplateIf, error) {
	return x.TemplateMap.Get(name)
}

func (x *XHandler) TemplateExecute(w io.Writer, name string, vars MapSA) error {
	tpl, err := x.TemplateMap.Get(name)
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

func (x *XHandler) TemplateExecuteSub(w io.Writer, name, subname string, vars MapSA) error {
	tpl, err := x.TemplateMap.Get(name)
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
	err := x.TemplateExecute(w, "error-page.gohtml", MapSA{
		"ErrorCode":   code,
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

func (x *XHandler) Json(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Info().Err(err).Msg("json encoding error")
		x.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

type MapSA = map[string]any
