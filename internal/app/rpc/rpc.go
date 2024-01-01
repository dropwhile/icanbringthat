package rpc

import (
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/crypto"
	"github.com/dropwhile/icbt/internal/mail"
	"github.com/dropwhile/icbt/internal/validate"
	"github.com/dropwhile/icbt/rpc/icbt"
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

type Options struct {
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

func (s *Server) GenHandler(prefix string) icbt.TwirpServer {
	twirpHandler := icbt.NewRpcServer(
		s,
		twirp.WithServerPathPrefix(prefix),
		twirp.WithServerHooks(
			&twirp.ServerHooks{
				RequestReceived: AuthHook(s.svc),
			},
		),
	)
	return twirpHandler
}
