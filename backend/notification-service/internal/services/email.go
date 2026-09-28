package services

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPEmailSender struct {
	host, port, username, password, from string
}

func NewSMTPEmailSender(host, port, username, password, from string) *SMTPEmailSender {
	return &SMTPEmailSender{host: host, port: port, username: username, password: password, from: from}
}

func (s *SMTPEmailSender) SendWelcome(ctx context.Context, recipient, name string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.host == "" || s.from == "" {
		return fmt.Errorf("email transport is not configured")
	}
	if strings.ContainsAny(recipient+s.from, "\r\n") {
		return fmt.Errorf("invalid email address")
	}
	if name == "" {
		name = "there"
	}
	body := fmt.Sprintf("To: %s\r\nSubject: Welcome\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\nHello %s,\r\n\r\nWelcome! Your account is ready.\r\n", recipient, name)
	var auth smtp.Auth
	if s.username != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}
	return smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{recipient}, []byte(body))
}
