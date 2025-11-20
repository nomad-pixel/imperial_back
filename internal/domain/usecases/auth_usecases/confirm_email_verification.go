package auth_usecases

import (
	"context"
	"errors"
	"time"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type ConfirmEmailVerificationUsecase interface {
	Execute(ctx context.Context, email, code string) error
}

type confirmEmailVerificationUsecase struct {
	verifyCodeRepo ports.VerifyCodeRepository
	userRepo       ports.UserRepository
}

func NewConfirmEmailVerificationUsecase(verifyCodeRepo ports.VerifyCodeRepository, userRepo ports.UserRepository) ConfirmEmailVerificationUsecase {
	return &confirmEmailVerificationUsecase{
		verifyCodeRepo: verifyCodeRepo,
		userRepo:       userRepo,
	}
}

func (u *confirmEmailVerificationUsecase) Execute(ctx context.Context, email, code string) error {
	verifyCode, err := u.verifyCodeRepo.GetVerifyCodeByEmailAndCodeAndType(
		ctx,
		email,
		code,
		entities.VerifyCodeTypeEmailVerification,
	)
	if err != nil {
		if errors.Is(err, apperrors.ErrVerifyCodeNotFound) {
			return apperrors.ErrVerifyCodeNotFound
		}
		return apperrors.Wrap(err, apperrors.ErrCodeDatabase, "Ошибка получения кода верификации")
	}

	if verifyCode.IsUsed {
		return apperrors.ErrVerifyCodeAlreadyUsed
	}

	now := time.Now()
	if verifyCode.ExpiresAt.Before(now) {
		return apperrors.ErrVerifyCodeExpired
	}

	verifyCode.IsUsed = true
	updatedCode, err := u.verifyCodeRepo.UpdateVerifyCode(ctx, verifyCode)
	if err != nil {
		if errors.Is(err, apperrors.ErrVerifyCodeNotFound) {
			return apperrors.Wrap(err, apperrors.ErrCodeDatabase, "Код верификации был удален во время обработки")
		}
		return apperrors.Wrap(err, apperrors.ErrCodeDatabase, "Ошибка обновления кода верификации")
	}

	_ = updatedCode

	_, err = u.userRepo.ConfirmEmailVerification(ctx, email)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeDatabase, "Ошибка подтверждения email пользователя")
	}

	return nil
}
