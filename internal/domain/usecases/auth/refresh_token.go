package usecases

import (
	"context"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type RefreshTokenUsecase interface {
	Execute(ctx context.Context, refreshToken string) (newAccessToken string, err error)
}

type refreshTokenUsecase struct {
	tokenSvc ports.TokenService
}

func NewRefreshTokenUsecase(tokenSvc ports.TokenService) RefreshTokenUsecase {
	return &refreshTokenUsecase{tokenSvc: tokenSvc}
}

func (u *refreshTokenUsecase) Execute(ctx context.Context, refreshToken string) (string, error) {
	return u.tokenSvc.RefreshAccessToken(refreshToken)
}
