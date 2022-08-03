package mailservice

import (
	"context"
	"github.com/decagonhq/meddle-api/config"
	"github.com/mailgun/mailgun-go/v4"
	"log"
	"os"
	"time"
)

type Mailer interface {
	SendSimpleMessage(UserEmail, EmailSubject, EmailBody string) error
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

func (mail Mailgun) SendSimpleMessage(UserEmail, EmailSubject, EmailBody string) error {
	EmailFrom := os.Getenv("MG_EMAIL_FROM")
	m := mail.Client.NewMessage(EmailFrom, EmailSubject, EmailBody, UserEmail)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, _, err := mail.Client.Send(ctx, m)
	if err != nil {
		log.Printf("could not send message %s", err)
	}
	return nil
}

