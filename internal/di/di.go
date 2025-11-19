package di

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/config"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/domain/usecases"
	token "github.com/nomad-pixel/imperial/internal/infrastructure/auth"
	"github.com/nomad-pixel/imperial/internal/infrastructure/email"
	"github.com/nomad-pixel/imperial/internal/infrastructure/repositories"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
)

func InitializeApp(ctx context.Context, dbURL string) (*App, error) {
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	userRepo := repositories.NewUserRepositoryImpl(db)
	verifyCodeRepo := repositories.NewVerifyCodeRepositoryImpl(db)

	emailConfig := config.LoadEmailConfig()
	fmt.Println("emailConfig", emailConfig)
	var emailService ports.EmailService

	if emailConfig.Provider == "smtp" {
		if emailConfig.SMTP.Username == "" {
			log.Println("⚠️  ВНИМАНИЕ: SMTP_USERNAME не установлен!")
		}
		if emailConfig.SMTP.Password == "" {
			log.Println("⚠️  ВНИМАНИЕ: SMTP_PASSWORD не установлен!")
		}

		smtpService, err := email.NewSMTPEmailService(email.SMTPConfig{
			Host:     emailConfig.SMTP.Host,
			Port:     emailConfig.SMTP.Port,
			Username: emailConfig.SMTP.Username,
			Password: emailConfig.SMTP.Password,
			From:     emailConfig.SMTP.From,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize SMTP service: %w", err)
		}
		emailService = smtpService
	} else {
		emailService = email.NewConsoleEmailService()
	}
	tokenSvc := token.NewJWTTokenService()

	signUpUsecase := usecases.NewSignUpUsecase(userRepo)
	sendEmailVerificationUsecase := usecases.NewSendEmailVerificationUsecase(
		userRepo,
		verifyCodeRepo,
		emailService,
	)
	confirmEmailVerificationUsecase := usecases.NewConfirmEmailVerificationUsecase(verifyCodeRepo, userRepo)
	signInUsecase := usecases.NewSignInUsecase(userRepo, tokenSvc)
	refreshTokenUsecase := usecases.NewRefreshTokenUsecase(tokenSvc)

	authHandler := auth.NewAuthHandler(
		signUpUsecase,
		sendEmailVerificationUsecase,
		confirmEmailVerificationUsecase,
		signInUsecase,
		refreshTokenUsecase,
	)

	app := NewApp(
		db,
		signUpUsecase,
		sendEmailVerificationUsecase,
		confirmEmailVerificationUsecase,
		signInUsecase,
		authHandler,
		tokenSvc,
	)

	return app, nil
}
