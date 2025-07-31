package models

import "github.com/go-mail/mail"

/*
	The mail.NewDialer(host, port, username, password) function creates a new SMTP dialer for sending emails.
	Each argument is required to establish a connection to your email provider’s SMTP server:

	host: The address of the SMTP server (e.g., "smtp.gmail.com"). This tells the dialer where to connect to send emails.
	port: The port number the SMTP server listens on (commonly 587 for TLS, 465 for SSL, or 25 for plain SMTP).
	username: The username for authenticating with the SMTP server (usually your email address).
	password: The password or app-specific password for authenticating the user.
	Summary:
	These arguments allow the dialer to connect securely to the SMTP server and authenticate your application, so it can send emails on your behalf.
*/

type EmailService struct {
	dialer *mail.Dialer
	to     string
}

func NewSMTPEmailService(host string, port int, username, password, to string) *EmailService {
	return &EmailService{
		dialer: mail.NewDialer(host, port, username, password),
		to:     to,
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

// SendEmailToRecipient sends an email to a specific recipient (overriding the default 'to' field)
func (es *EmailService) SendEmailToRecipient(from, to, subject, body string) error {
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	return es.dialer.DialAndSend(m)
}
