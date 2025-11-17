package ports

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
)

type VerifyCodeRepository interface {
	CreateVerifyCode(ctx context.Context, code string, userID int64, verifyCodeType entities.VerifyCodeType) (*entities.VerifyCode, error)
	GetVerifyCodeByUserIDAndType(ctx context.Context, userID int64, verifyCodeType entities.VerifyCodeType) (*entities.VerifyCode, error)
	UpdateVerifyCode(ctx context.Context, code string, userID int64, verifyCodeType entities.VerifyCodeType) (*entities.VerifyCode, error)
}
