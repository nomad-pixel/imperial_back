package usecases

import (
	"context"
	"time"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
	"github.com/nomad-pixel/imperial/pkg/utils"
)

type SendEmailVerificationUsecase interface {
	Execute(ctx context.Context, email string) error
}

type sendEmailVerificationUsecase struct {
	userRepo       ports.UserRepository
	verifyCodeRepo ports.VerifyCodeRepository
	emailService   ports.EmailService
}

func NewSendEmailVerificationUsecase(
	userRepo ports.UserRepository,
	verifyCodeRepo ports.VerifyCodeRepository,
	emailService ports.EmailService,
) SendEmailVerificationUsecase {
	return &sendEmailVerificationUsecase{
		userRepo:       userRepo,
		verifyCodeRepo: verifyCodeRepo,
		emailService:   emailService,
	}
}

func (u *sendEmailVerificationUsecase) Execute(ctx context.Context, email string) error {
	user, err := u.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user.VerifiedAt {
		return apperrors.New(apperrors.ErrCodeValidation, "Email уже верифицирован")
	}

	code, err := utils.GenerateVerificationCode(6)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeInternal, "Ошибка генерации кода верификации")
	}

	existingCode, err := u.verifyCodeRepo.GetVerifyCodeByUserIDAndType(ctx, user.ID, entities.VerifyCodeTypeEmailVerification)

	if existingCode != nil {
		existingCode.ExpiresAt = time.Now().Add(5 * time.Minute)
		existingCode.Code = code
		existingCode.IsUsed = false
		_, err = u.verifyCodeRepo.UpdateVerifyCode(ctx, existingCode)
		if err != nil {
			return err
		}
	}

	_, err = u.verifyCodeRepo.CreateVerifyCode(
		ctx,
		code,
		user.ID,
		entities.VerifyCodeTypeEmailVerification,
		time.Now().Add(5*time.Minute),
	)
	if err != nil {
		return err
	}

	err = u.emailService.SendVerificationCode(ctx, email, code)
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeExternal, "Ошибка отправки email")
	}

	return nil
}
