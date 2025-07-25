package models

import "github.com/go-mail/mail"

type EmailService struct {
	dialer *mail.Dialer
	to string
}

func NewSMTPEmailService(host string, port int, username, password, to string) *EmailService {
	return &EmailService{
		dialer: mail.NewDialer(host, port, username, password),
		to: to,
	}
}

func (es *EmailService) SendEmail(from, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", es.to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return es.dialer.DialAndSend(m)
}