package handler

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/redis/go-redis/v9"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/crypto"
	"github.com/dropwhile/icbt/internal/logger"
	"github.com/dropwhile/icbt/internal/mail"
	"github.com/dropwhile/icbt/internal/session"
)

type Handler struct {
	Db        model.PgxHandle
	Redis     *redis.Client
	Templates resources.TGetter
	SessMgr   *session.SessionMgr
	Mailer    mail.MailSender
	MAC       *crypto.MAC
	Service   service.Servicer
	BaseURL   string
	IsProd    bool
}

func NewHandler(
	db model.PgxHandle, redis *redis.Client, templates resources.TGetter,
	sessmgr *session.SessionMgr, mailer mail.MailSender, mac *crypto.MAC,
	baseURL string, isProd bool,
) *Handler {
	handler := &Handler{
		Db:        db,
		Redis:     redis,
		Templates: templates,
		SessMgr:   sessmgr,
		Mailer:    mailer,
		MAC:       mac,
		Service:   &service.Service{Db: db},
		BaseURL:   baseURL,
		IsProd:    isProd,
	}
	return handler
}

func (x *Handler) Template(name string) (resources.TemplateIf, error) {
	return x.Templates.Get(name)
}

func (x *Handler) TemplateExecute(w io.Writer, name string, vars MapSA) error {
	tpl, err := x.Templates.Get(name)
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
	tpl, err := x.Templates.Get(name)
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
