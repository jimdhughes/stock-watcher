package util

import (
	"errors"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

const (
	SMTP_HOST     = "SMTP_HOST"
	SMPT_PORT     = "SMTP_PORT"
	SMTP_EMAIL    = "SMTP_EMAIL"
	SMTP_PASSWORD = "SMTP_PASSWORD"
)

type Mailer struct {
	SmtpHost string
	Password string
	Email    string
	SmtpPort string
}

var AppMailer *Mailer

func (m *Mailer) SendMail(to []string, message, subject string) error {
	if m.SmtpHost == "" || m.Email == "" || m.Password == "" || m.SmtpPort == "" {
		return errors.New("Mailer is not initialized")
	}
	auth := smtp.PlainAuth("", m.Email, m.Password, m.SmtpHost)
	msg := m.formatMessage(subject, message, to)
	err := smtp.SendMail(fmt.Sprintf("%s:%s", m.SmtpHost, m.SmtpPort), auth, m.Email, to, []byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (m *Mailer) formatMessage(subject, message string, to []string) string {
	str := fmt.Sprintf("to: %s\r\nSubject: %s\r\n\r\n%s", to, subject, message)
	log.Println(str)
	return str
}

func init() {
	AppMailer = &Mailer{
		SmtpHost: os.Getenv(SMTP_HOST),
		SmtpPort: os.Getenv(SMPT_PORT),
		Email:    os.Getenv(SMTP_EMAIL),
		Password: os.Getenv(SMTP_PASSWORD),
	}
}
