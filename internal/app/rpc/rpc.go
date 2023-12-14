package rpc

import (
	"github.com/redis/go-redis/v9"
	"github.com/twitchtv/twirp"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/app/resources"
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
	BaseURL   string
	IsProd    bool
}

func (s *Server) GenHandler(prefix string) icbt.TwirpServer {
	twirpHandler := icbt.NewRpcServer(
		s,
		twirp.WithServerPathPrefix(prefix),
		twirp.WithServerHooks(
			&twirp.ServerHooks{
				RequestReceived: AuthHook(s.Db),
			},
		),
	)
	return twirpHandler
}
