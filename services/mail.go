package services

import (
	"github.com/decagonhq/meddle-api/config"
	"github.com/mailgun/mailgun-go/v4"
)

type Mailer interface {
	SendMail(toEmail, title, body, template string, values map[string]interface{}) error
}
type Mailgun struct {
	Client *mailgun.MailgunImpl
	Conf   *config.Config
}

// NewMailService instantiates a mail service
func NewMailService(conf *config.Config) Mailer {
	domain := conf.MgDomain
	apiKey := conf.MailgunApiKey
	return &Mailgun{
		Client: mailgun.NewMailgun(domain, apiKey),
		Conf:   conf,
	}
}

func (mail *Mailgun) SendMail(toEmail, title, body, link string, v map[string]interface{}) error {
	return nil
}
