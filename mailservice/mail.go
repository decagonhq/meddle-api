package mailservice

import (
	"context"
	"errors"
	"github.com/decagonhq/meddle-api/config"
	"github.com/mailgun/mailgun-go/v4"
	"os"
	"time"
)

type Mailer interface {
	SendVerifyAccount(userEmail, link string) (string, error)
}
type Mailgun struct {
	Client *mailgun.MailgunImpl
	Conf *config.Config
}

// NewMailService instantiates a mail service
func NewMailService(mail *mailgun.MailgunImpl, conf *config.Config) Mailer {
	return &Mailgun{
		Client: mail,
		Conf:    conf,
	}
}

func StartMailGun() *mailgun.MailgunImpl {
	domain := os.Getenv("MG_DOMAIN")
	apiKey := os.Getenv("MG_PUBLIC_API_KEY")
	mail := mailgun.NewMailgun(domain, apiKey)
	return mail
}

func (mail *Mailgun) SendVerifyAccount(userEmail, link string) (string, error) {
	EmailFrom := os.Getenv("MG_EMAIL_FROM")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	m := mail.Client.NewMessage(EmailFrom, "Verify Account", "")
	m.SetTemplate("verify.account")
	if err := m.AddRecipient(userEmail); err != nil {
		return "", errors.New("could not add recipient")
	}
	err := m.AddVariable("link", link)
	if err != nil {
		return "", err
	}
	_, _, err = mail.Client.Send(ctx, m)
	if err != nil {
		return "", err
	}
	return "", err
}
