package domain

import "context"

type Email struct {
	To      string
	Subject string
	Body    string
}

type EmailSender interface {
	Send(ctx context.Context, email *Email) error
	SendVerificationEmail(ctx context.Context, email, token string) error
}
