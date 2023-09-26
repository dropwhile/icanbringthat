package util

import (
	"fmt"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/dropwhile/refid"
	"github.com/rs/zerolog/log"
)

type Mail struct {
	Sender    string
	To        []string
	Subject   string
	BodyPlain string
	BodyHtml  string
}

type Mailer struct {
	hostname string
	hostPort string
	user     string
	auth     smtp.Auth
}

func (m *Mailer) SendRaw(mail *Mail) error {
	if mail.BodyHtml == "" && mail.BodyPlain == "" {
		return fmt.Errorf("no content")
	}
	var buf strings.Builder
	boundary := refid.Must(refid.New())
	// write headers, set up boundary
	buf.WriteString(fmt.Sprintf("From: %s\r\n", mail.Sender))
	buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";")))
	buf.WriteString(fmt.Sprintf("Subject: %s\r\n", mail.Subject))
	buf.WriteString("MIME-Version: 1.0\r\n")
	buf.WriteString(fmt.Sprintf(
		"Content-Type: multipart/alternative;boundary=%s\r\n", boundary,
	))
	// add body contents in boundary delineated sections
	if mail.BodyPlain != "" {
		buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
		buf.WriteString("Content-type: text/plain;charset=utf-8\r\n")
		buf.WriteString(fmt.Sprintf("\r\n%s\r\n", mail.BodyPlain))
	}
	if mail.BodyHtml != "" {
		buf.WriteString(fmt.Sprintf("\r\n--%s\r\n", boundary))
		buf.WriteString("Content-type: text/html;charset=utf-8\r\n")
		buf.WriteString(fmt.Sprintf("\r\n%s\r\n", mail.BodyHtml))
	}
	buf.WriteString(fmt.Sprintf("\r\n--%s--", boundary))

	message := buf.String()
	log.Debug().
		Str("message", message).
		Msg("sending email")
	err := smtp.SendMail(m.hostPort, m.auth, mail.Sender, mail.To, []byte(message))
	if err != nil {
		log.Info().Err(err).Msg("error sending email")
	}
	return err
}

func (m *Mailer) Send(from string, to []string, subject, bodyPlain, bodyHtml string) error {
	if from == "" {
		from = m.user
	}
	mail := &Mail{
		Sender:    from,
		To:        to,
		Subject:   subject,
		BodyPlain: bodyPlain,
		BodyHtml:  bodyHtml,
	}
	return m.SendRaw(mail)
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
