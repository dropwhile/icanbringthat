// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/resources"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/crypto"
	"github.com/dropwhile/icanbringthat/internal/logger"
	"github.com/dropwhile/icanbringthat/internal/mail"
	"github.com/dropwhile/icanbringthat/internal/session"
	"github.com/dropwhile/icanbringthat/internal/validate"
)

type Handler struct {
	redis     *redis.Client
	templates resources.TGetter
	sessMgr   session.SessionManager
	mailer    mail.MailSender
	cMAC      crypto.HMACer
	svc       service.Servicer
	baseURL   string
	isProd    bool
}

type Options struct { // betteralign:ignore
	Db           model.PgxHandle        `validate:"required"`
	Redis        *redis.Client          `validate:"required"`
	Templates    resources.TGetter      `validate:"required"`
	SessMgr      session.SessionManager `validate:"required"`
	Mailer       mail.MailSender        `validate:"required"`
	HMACKeyBytes []byte                 `validate:"required"`
	BaseURL      string                 `validate:"required"`
	IsProd       bool
}

func New(opts Options) (*Handler, error) {
	err := validate.Validate.Struct(opts)
	if err != nil {
		badField := validate.GetErrorField(err)
		slog.
			With("field", badField).
			With("error", err).
			Info("bad field value")
		return nil, fmt.Errorf("bad options value: %s", badField)
	}

	cMAC := crypto.NewMAC(opts.HMACKeyBytes)
	handler := &Handler{
		redis:     opts.Redis,
		templates: opts.Templates,
		sessMgr:   opts.SessMgr,
		mailer:    opts.Mailer,
		cMAC:      cMAC,
		baseURL:   opts.BaseURL,
		isProd:    opts.IsProd,
		svc:       &service.Service{Db: opts.Db},
	}
	return handler, nil
}

func (x *Handler) Template(name string) (resources.TExecuter, error) {
	return x.templates.Get(name)
}

func (x *Handler) TemplateExecute(w io.Writer, name string, vars MapSA) error {
	tpl, err := x.templates.Get(name)
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelInfo,
			context.Background(),
			"template locate error",
			"error", err,
			"tpl", name)
		return err
	}
	err = tpl.Execute(w, vars)
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelInfo,
			context.Background(),
			"template execute error",
			"error", err, "tpl", name, "vars", vars)
		return err
	}
	return nil
}

func (x *Handler) TemplateExecuteSub(w io.Writer, name, subname string, vars MapSA) error {
	tpl, err := x.templates.Get(name)
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelInfo,
			context.Background(),
			"template locate error",
			"error", err, "tpl", name, "sub", subname)
		return err
	}
	err = tpl.ExecuteTemplate(w, subname, vars)
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelInfo,
			context.Background(),
			"template execute error",
			"error", err, "tpl", name, "sub", subname, "vars", vars)
		return err
	}
	return nil
}

func (x *Handler) Json(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		logger.LogSkip(slog.Default(), 1, slog.LevelInfo,
			context.Background(),
			"json encoding error",
			"error", err)
		x.Error(w, "encoding error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, _ = w.Write(response)
}

type MapSA = map[string]any
