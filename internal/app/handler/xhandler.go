package handler

import (
	"encoding/json"
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

type Handler struct {
	Db          model.PgxHandle
	Redis       *redis.Client
	TemplateMap resources.TemplateMap
	SessMgr     *session.SessionMgr
	Mailer      mail.MailSender
	MAC         *crypto.MAC
	BaseURL     string
	IsProd      bool
}

func (x *Handler) Template(name string) (resources.TemplateIf, error) {
	return x.TemplateMap.Get(name)
}

func (x *Handler) TemplateExecute(w io.Writer, name string, vars MapSA) error {
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

func (x *Handler) TemplateExecuteSub(w io.Writer, name, subname string, vars MapSA) error {
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

func (x *Handler) Json(w http.ResponseWriter, code int, payload interface{}) {
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
