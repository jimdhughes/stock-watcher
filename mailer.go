package main

import (
	"errors"
	"fmt"
	"net/smtp"
)

type Mailer struct {
	SmtpHost string
	Password string
	Email string
	SmtpPort string
}

func (m *Mailer) SendMail(to []string, message, subject string) error {
	if (m.SmtpHost == "" || m.Email == "" || m.Password == "" || m.SmtpPort == "") {
		return errors.New("Mailer is not initialized")
	}
	auth := smtp.PlainAuth("", m.Email, m.Password, m.SmtpHost)
	msg := m.formatMessage(subject, message)
	err := smtp.SendMail(m.SmtpHost+":"+m.SmtpPort,auth, m.Email, to, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (m *Mailer) formatMessage(subject, message string) string {
	return fmt.Sprintf("Subject: %s\r\n\r\n%s\n", subject, message)
}