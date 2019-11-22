package mail

import (
	"crypto/tls"
	"fmt"
	"github.com/FlagField/FlagField-Server/internal/pkg/config"
	"github.com/mukeran/email"
	"net/smtp"
	"net/textproto"
	"time"
)

type Mail struct {
	config config.MailConfig
	pool   *email.Pool
}

func (m *Mail) Mail(to []string, subject string, text []byte, html []byte) *email.Email {
	return &email.Email{
		From:    fmt.Sprintf("%s <%s>", m.config.SenderName, m.config.Address),
		To:      to,
		Subject: subject,
		Text:    text,
		HTML:    html,
		Headers: textproto.MIMEHeader{},
	}
}

func (m *Mail) Send(e *email.Email) error {
	return m.pool.Send(e, 3*time.Second)
	//return e.Send(fmt.Sprintf("%s:%d", m.config.Host, m.config.Port), smtp.PlainAuth("", m.config.Username, m.config.Password, m.config.Host))
}

func (m *Mail) Test(to string) error {
	e := m.Mail([]string{to}, "Test Mail", []byte("This is a test mail from a FlagField-Server which indicates that the SMTP settings work!"), nil)
	return m.Send(e)
}

func New(config config.MailConfig) (*Mail, error) {
	var (
		p   *email.Pool
		err error
	)
	if !config.UseTLS {
		p, err = email.NewPool(fmt.Sprintf("%s:%d", config.Host, config.Port), 1, smtp.PlainAuth("", config.Username, config.Password, config.Host), false)
	} else {
		p, err = email.NewPool(fmt.Sprintf("%s:%d", config.Host, config.Port), 1, smtp.PlainAuth("", config.Username, config.Password, config.Host), true, &tls.Config{InsecureSkipVerify: true, ServerName: config.Host})
	}
	if err != nil {
		return nil, err
	}
	return &Mail{config: config, pool: p}, nil
}
