package ports

import "context"

type EmailService interface {
	SendVerificationCode(ctx context.Context, email, code string) error
	SendPasswordResetCode(ctx context.Context, email, code string) error
}
