package _rdPartyService

import (
	"github.com/mailgun/mailgun-go"
	"os"
)

var (
	mailgunDomain = os.Getenv("MAILGUN_DOMAIN")
	mailgunApiKey = os.Getenv("MAILGUN_API_KEY")
)

func SendMail(message, receiver string) (string, error) {
	mg := mailgun.NewMailgun(mailgunDomain, mailgunApiKey)
	m := mg.NewMessage(
		"Authenticator <mailgun@"+mailgunDomain+">",
		"Identity Authentication Results",
		message,
		receiver,
	)
	_, id, err := mg.Send(m)
	return id, err
}
