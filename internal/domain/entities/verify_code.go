package entities

import "time"

type VerifyCodeType string

const (
	VerifyCodeTypeEmailVerification VerifyCodeType = "email_verification"
	VerifyCodeTypePasswordReset     VerifyCodeType = "password_reset"
)

type VerifyCode struct {
	ID        int64
	UserID    int64
	Code      string
	Type      VerifyCodeType
	IsUsed    bool
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
