package ports

import (
	"context"
	"time"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type VerifyCodeRepository interface {
	CreateVerifyCode(ctx context.Context, code string, userID int64, verifyCodeType entities.VerifyCodeType, expiresAt time.Time) (*entities.VerifyCode, error)
	GetVerifyCodeByUserIDAndType(ctx context.Context, userID int64, verifyCodeType entities.VerifyCodeType) (*entities.VerifyCode, error)
	GetVerifyCodeByEmailAndCodeAndType(ctx context.Context, email, code string, verifyCodeType entities.VerifyCodeType) (*entities.VerifyCode, error)
	GetVerifyCodeByCodeAndType(ctx context.Context, code string, verifyCodeType entities.VerifyCodeType) (*entities.VerifyCode, error)
	UpdateVerifyCode(ctx context.Context, verifyCode *entities.VerifyCode) (*entities.VerifyCode, error)
}
