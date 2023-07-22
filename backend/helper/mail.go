package helper

import (
	"strconv"

	"gopkg.in/gomail.v2"
)

// Config for send email message to client
type MailerConfig struct {
	// SMTP Host
	Host string

	// SMTP Port
	Port string

	// Name to be displayed
	SenderName string

	// Email
	Email string

	// Password email
	Password string
}

func NewMailer(host, port, senderName, email, password string) *MailerConfig {
	return &MailerConfig{
		Host:       host,
		Port:       port,
		SenderName: senderName,
		Email:      email,
		Password:   password,
	}
}

func (mailConfig *MailerConfig) SendMail(to, subject, message string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", mailConfig.SenderName)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", message)

	port, _ := strconv.Atoi(mailConfig.Port)

	dialer := gomail.NewDialer(
		mailConfig.Host,
		port,
		mailConfig.Email,
		mailConfig.Password,
	)

	err := dialer.DialAndSend(mailer)

	if err != nil {
		return err
	}

	return nil
}
