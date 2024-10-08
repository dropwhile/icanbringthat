// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package rpc

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"

	"connectrpc.com/connect"
	connectValidate "connectrpc.com/validate"
	"github.com/redis/go-redis/v9"

	"github.com/dropwhile/icanbringthat/internal/app/model"
	"github.com/dropwhile/icanbringthat/internal/app/resources"
	"github.com/dropwhile/icanbringthat/internal/app/service"
	"github.com/dropwhile/icanbringthat/internal/crypto"
	"github.com/dropwhile/icanbringthat/internal/mail"
	"github.com/dropwhile/icanbringthat/internal/validate"
	"github.com/dropwhile/icanbringthat/rpc/icbt/rpc/v1/rpcv1connect"
)

type Server struct {
	redis     *redis.Client
	templates resources.TGetter
	mailer    mail.MailSender
	cMAC      crypto.HMACer
	svc       service.Servicer
	baseURL   string
	isProd    bool
}

type Options struct { // betteralign:ignore
	Db           model.PgxHandle   `validate:"required"`
	Redis        *redis.Client     `validate:"required"`
	Templates    resources.TGetter `validate:"required"`
	Mailer       mail.MailSender   `validate:"required"`
	HMACKeyBytes []byte            `validate:"required"`
	BaseURL      string            `validate:"required"`
	IsProd       bool
}

func New(opts Options) (*Server, error) {
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
	svr := &Server{
		redis:     opts.Redis,
		templates: opts.Templates,
		mailer:    opts.Mailer,
		cMAC:      cMAC,
		svc: &service.Service{
			Db: opts.Db,
		},
		baseURL: opts.BaseURL,
		isProd:  opts.IsProd,
	}
	return svr, nil
}

func (s *Server) GenHandler() http.Handler {
	rpcValidateInterceptor, err := connectValidate.NewInterceptor()
	if err != nil {
		log.Fatal(err)
	}

	interceptors := connect.WithInterceptors(
		// NewAuthInterceptor(s.svc),
		rpcValidateInterceptor,
	)
	api := http.NewServeMux()
	api.Handle(rpcv1connect.NewIcbtRpcServiceHandler(s, interceptors))
	return api
}
