package services

import (

	"context"
	"errors"
	"github.com/decagonhq/meddle-api/config"
	"github.com/mailgun/mailgun-go/v4"
	"time"
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


func (m *Mailgun) SendMail(toEmail, subject, body, template string, values map[string]interface{}) error {
	message := m.Client.NewMessage(m.Conf.EmailFrom, subject, body)
	message.SetTemplate(template)
	if err := message.AddRecipient(toEmail); err != nil {
		return errors.New("could not add recipient")
	}
	for k, v := range values {
		err := message.AddVariable(k, v)
		if err != nil {
			return err
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, _, err := m.Client.Send(ctx, message)
	return err
}
