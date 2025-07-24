package models

import "github.com/go-mail/mail"

type EmailService struct {
	dialer *mail.Dialer
	from string
}

func NewSMTPEmailService(host string, port int, username, password, from string) *EmailService {
	return &EmailService{
		dialer: mail.NewDialer(host, port, username, password),
		from: from,
	}
}

func (es *EmailService) SendEmail(to, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", es.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return es.dialer.DialAndSend(m)
}