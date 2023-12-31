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
	redis     *redis.Client
	templates resources.TGetter
	mailer    mail.MailSender
	cMAC      crypto.HMACer
	service   service.Servicer
	baseURL   string
	isProd    bool
}

type Options struct {
	Db           model.PgxHandle
	Redis        *redis.Client
	Templates    resources.TGetter
	Mailer       mail.MailSender
	HMACKeyBytes []byte
	BaseURL      string
	IsProd       bool
}

func New(opts Options) *Server {
	cMAC := crypto.NewMAC(opts.HMACKeyBytes)
	svr := &Server{
		redis:     opts.Redis,
		templates: opts.Templates,
		mailer:    opts.Mailer,
		cMAC:      cMAC,
		service: &service.Service{
			Db: opts.Db,
		},
		baseURL: opts.BaseURL,
		isProd:  opts.IsProd,
	}
	return svr
}

func (s *Server) GenHandler(prefix string) icbt.TwirpServer {
	twirpHandler := icbt.NewRpcServer(
		s,
		twirp.WithServerPathPrefix(prefix),
		twirp.WithServerHooks(
			&twirp.ServerHooks{
				RequestReceived: AuthHook(s.service),
			},
		),
	)
	return twirpHandler
}
