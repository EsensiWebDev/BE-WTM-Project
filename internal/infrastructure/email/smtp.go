package email

import (
	"context"
	"fmt"
	"gopkg.in/gomail.v2"
	"wtm-backend/config"
	"wtm-backend/pkg/constant"
	"wtm-backend/pkg/logger"
)

type SMTPEmailSender struct {
	Host          string
	Port          int
	Username      string
	Password      string
	From          string
	SupportEmail  string
	EmailProvider string
}

func NewSMTPEmailSender(config *config.Config) *SMTPEmailSender {
	return &SMTPEmailSender{
		Host:          config.EmailHost,
		Port:          config.EmailPort,
		Username:      config.EmailUser,
		Password:      config.EmailPass,
		From:          config.EmailFrom,
		SupportEmail:  config.SupportEmail,
		EmailProvider: config.EmailProvider,
	}
}

func (s *SMTPEmailSender) Send(ctx context.Context, to, subject, bodyHTML, bodyText string) error {
	if to == "" {
		return fmt.Errorf("email recipient is empty")
	}

	if to == constant.SupportEmail {
		to = s.SupportEmail
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", bodyHTML)
	m.AddAlternative("text/plain", bodyText)

	d := gomail.NewDialer(s.Host, s.Port, s.Username, s.Password)

	if s.EmailProvider == "mailhog" {
		d.Auth = nil
	}
	if s.EmailProvider == "ses" {
		m.SetHeader("Reply-To", s.SupportEmail)
		m.SetHeader("Return-Path", s.SupportEmail)
	}
	if s.EmailProvider == "gmail" && s.Username == "" {
		logger.Warn(ctx, "Gmail config missing, fallback to MailHog")
		s.EmailProvider = "mailhog"
		d.Auth = nil
	}

	logger.Info(ctx, fmt.Sprintf("Sending email via provider: %s", s.EmailProvider))

	err := d.DialAndSend(m)
	if err != nil {
		logger.Error(ctx, fmt.Sprintf("Failed to send email: %v", err.Error()))
		return err
	}

	return nil
}
