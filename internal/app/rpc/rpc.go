package rpc

import (
	"github.com/redis/go-redis/v9"

	"github.com/dropwhile/icbt/internal/app/model"
	"github.com/dropwhile/icbt/internal/crypto"
	"github.com/dropwhile/icbt/internal/mail"
	"github.com/dropwhile/icbt/resources"
)

type Server struct {
	Db          model.PgxHandle
	Redis       *redis.Client
	TemplateMap resources.TemplateMap
	Mailer      mail.MailSender
	MAC         *crypto.MAC
	BaseURL     string
	IsProd      bool
}
