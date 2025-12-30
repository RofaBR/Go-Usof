package services

import (
	"context"
	"fmt"
	"os"

	"github.com/RofaBR/Go-Usof/internal/config"
	"github.com/RofaBR/Go-Usof/internal/domain"
	"gopkg.in/gomail.v2"
)

type SMPTSender struct {
	config config.SenderConfig
}

func NewSMPTSender(config config.SenderConfig) *SMPTSender {
	return &SMPTSender{config: config}
}
func (s *SMPTSender) Send(ctx context.Context, email *domain.Email) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.config.FromEmail)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/html", email.Body)

	d := gomail.NewDialer(
		s.config.SMTPHost,
		s.config.SMTPPort,
		s.config.FromEmail,
		s.config.Password,
	)

	err := d.DialAndSend(m)
	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", email.To, err)
	}

	return nil
}

func (s *SMPTSender) SendVerificationEmail(ctx context.Context, email, token string) error {
	baseUrl := os.Getenv("BASE_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:8080/api"
	}
	verifyURl := fmt.Sprintf("%s/auth/verify?token=%s", baseUrl, token)
	msg := domain.Email{
		To:      email,
		Subject: "Verify your email",
		Body:    fmt.Sprintf("<a href='%s'>Click here to verify your email</a>", verifyURl),
	}
	return s.Send(ctx, &msg)
}
