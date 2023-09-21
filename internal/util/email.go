package util

import (
	"net/smtp"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
)

type Mailer struct {
	hostname string
	hostPort string
	user     string
	auth     smtp.Auth
}

func (m *Mailer) Send(from, to, message string) error {
	if from == "" {
		from = m.user
	}
	log.Debug().
		Str("to", to).
		Str("message", message).
		Msg("sending email")
	err := smtp.SendMail(m.hostPort, m.auth, from, []string{to}, []byte(message))
	if err != nil {
		log.Info().Err(err).Msg("error sending email")
	}
	return err
}

func NewMailer(host string, port int, hostname string, user, pass string) *Mailer {
	auth := smtp.PlainAuth("", user, pass, hostname)
	return &Mailer{
		hostname: hostname,
		hostPort: strings.Join([]string{host, strconv.Itoa(port)}, ":"),
		user:     user,
		auth:     auth,
	}
}
