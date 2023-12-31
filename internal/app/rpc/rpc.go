package rpc

import (
	"github.com/redis/go-redis/v9"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
	"github.com/dropwhile/icbt/internal/app/service"
	"github.com/dropwhile/icbt/internal/crypto"
	"github.com/dropwhile/icbt/internal/mail"
	"github.com/dropwhile/icbt/rpc/icbt"
)

type Server struct {
	Db        model.PgxHandle
	Redis     *redis.Client
	Templates resources.TGetter
	Mailer    mail.MailSender
	MAC       *crypto.MAC
	Service   service.Servicer
	BaseURL   string
	IsProd    bool
}

func (s *Server) GenHandler(prefix string) icbt.TwirpServer {
	twirpHandler := icbt.NewRpcServer(
		s,
		twirp.WithServerPathPrefix(prefix),
		twirp.WithServerHooks(
			&twirp.ServerHooks{
				RequestReceived: AuthHook(s.Service),
			},
		),
	)
	return twirpHandler
}

func NewServer(
	db model.PgxHandle, redis *redis.Client, templates resources.TGetter,
	mailer mail.MailSender, mac *crypto.MAC, baseURL string, isProd bool,
) *Server {
	svr := &Server{
		Db:        db,
		Redis:     redis,
		Templates: templates,
		Mailer:    mailer,
		MAC:       mac,
		Service: &service.Service{
			Db: db,
		},
		BaseURL: baseURL,
		IsProd:  isProd,
	}
	return svr
}
