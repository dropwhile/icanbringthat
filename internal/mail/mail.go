// Copyright (c) 2024 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package mail

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/dropwhile/refid/v2"

	"github.com/dropwhile/icanbringthat/internal/logger"
	"github.com/dropwhile/icanbringthat/internal/util"
)

type MailHeader map[string]string

type Mail struct {
	ExtraHeaders MailHeader
	Sender       string
	Subject      string
	BodyPlain    string
	BodyHtml     string
	To           []string
}

type Mailer struct {
	hostname    string
	hostPort    string
	user        string
	auth        smtp.Auth
	defaultFrom string
}

//go:generate tool mockgen -source mail.go -destination mockmail/mailsender_mock.go -package mockmail
type MailSender interface {
	SendRaw(*Mail) error
	Send(string, []string, string, string, string, MailHeader) error
	SendAsync(string, []string, string, string, string, MailHeader)
}

func (m *Mailer) SendRaw(mail *Mail) error {
	if mail.BodyHtml == "" && mail.BodyPlain == "" {
		return fmt.Errorf("no content")
	}
	var buf strings.Builder
	boundary := util.Must(refid.New())
	// write headers, set up boundary
	for k, v := range mail.ExtraHeaders {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	buf.WriteString(fmt.Sprintf("From: ICanBringThat <%s>\r\n", mail.Sender))
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
	slog.Debug("sending email", slog.String("message", message))
	err := smtp.SendMail(m.hostPort, m.auth, mail.Sender, mail.To, []byte(message))
	if err != nil {
		slog.Error("error sending email", logger.Err(err))
	}
	return err
}

func (m *Mailer) Send(from string, to []string, subject, bodyPlain, bodyHtml string, extraHeaders MailHeader) error {
	if from == "" {
		from = m.defaultFrom
	}
	mail := &Mail{
		Sender:       from,
		To:           to,
		Subject:      subject,
		BodyPlain:    bodyPlain,
		BodyHtml:     bodyHtml,
		ExtraHeaders: extraHeaders,
	}
	return m.SendRaw(mail)
}

func (m *Mailer) SendAsync(from string, to []string, subject, bodyPlain, bodyHtml string, extraHeaders MailHeader) {
	go func() {
		err := m.Send(from, to, subject, bodyPlain, bodyHtml, extraHeaders)
		if err != nil {
			slog.Error("error sending email", logger.Err(err))
		}
	}()
}

func NewMailer(conf *Config) *Mailer {
	auth := smtp.PlainAuth("", conf.User, conf.Pass, conf.Hostname)
	return &Mailer{
		hostname:    conf.Hostname,
		hostPort:    strings.Join([]string{conf.Host, strconv.Itoa(conf.Port)}, ":"),
		user:        conf.User,
		auth:        auth,
		defaultFrom: conf.DefaultFrom,
	}
}
