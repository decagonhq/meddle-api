package mailservice

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"os"
	"time"
)

type Mailgun struct {
	Client *mailgun.MailgunImpl
}

func (mail *Mailgun) Init() {
	domain := os.Getenv("MG_DOMAIN")
	apiKey := os.Getenv("MG_PUBLIC_API_KEY")
	mail.Client = mailgun.NewMailgun(domain, apiKey)
}
func (mail Mailgun) SendSimpleMessage(UserEmail, EmailSubject, EmailBody string) (string, error) {
	EmailFrom := os.Getenv("MG_EMAIL_FROM")

	m := mail.Client.NewMessage(EmailFrom, EmailSubject, EmailBody, UserEmail)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	//m := mail.Client.NewMessage(EmailFrom, "Verify Account", "")
	res, _, err := mail.Client.Send(ctx, m)
	if err != nil {
		return "", err
	}
	return res, nil
}

