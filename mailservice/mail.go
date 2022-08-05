package mailservice

import (
	"context"
	"errors"
	"github.com/decagonhq/meddle-api/config"
	"github.com/mailgun/mailgun-go/v4"
	"time"
)

type Mailer interface {
	SendVerifyAccount(userEmail, link string) error
	SendResetPassword(userEmail, link string) (string, error)
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

func (mail *Mailgun) SendVerifyAccount(userEmail, link string) error {
	EmailFrom := mail.Conf.EmailFrom

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	m := mail.Client.NewMessage(EmailFrom, "Verify Account", "")
	m.SetTemplate("verify.account")
	if err := m.AddRecipient(userEmail); err != nil {
		return errors.New("could not add recipient")
	}
	err := m.AddVariable("link", link)
	if err != nil {
		return err
	}
	_, _, err = mail.Client.Send(ctx, m)
	return err
}

func (mail *Mailgun) SendResetPassword(userEmail, link string) (string, error) {
	EmailFrom := mail.Conf.EmailFrom

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	m := mail.Client.NewMessage(EmailFrom, "Reset Password", "")
	m.SetTemplate("reset.password")
	if err := m.AddRecipient(userEmail); err != nil {
		return "", err
	}

	err := m.AddVariable("link", link)
	if err != nil {
		return "", err
	}

	res, _, errr := mail.Client.Send(ctx, m)
	if errr != nil {
		return "", errr
	}
	return res, nil
}
