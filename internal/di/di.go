package di

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/config"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/internal/domain/usecases"
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
	var emailService ports.EmailService

	if emailConfig.Provider == "smtp" {
		log.Println("📧 Используется SMTP провайдер для отправки email")
		log.Printf("📧 SMTP Host: %s:%s", emailConfig.SMTP.Host, emailConfig.SMTP.Port)
		log.Printf("📧 SMTP From: %s", emailConfig.SMTP.From)
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
		log.Println("📧 Используется Console провайдер (email в консоль)")
		log.Println("💡 Для реальной отправки email установите EMAIL_PROVIDER=smtp")
		emailService = email.NewConsoleEmailService()
	}

	signUpUsecase := usecases.NewSignUpUsecase(userRepo)
	sendEmailVerificationUsecase := usecases.NewSendEmailVerificationUsecase(
		userRepo,
		verifyCodeRepo,
		emailService,
	)

	authHandler := auth.NewAuthHandler(signUpUsecase, sendEmailVerificationUsecase)

	app := NewApp(db, signUpUsecase, sendEmailVerificationUsecase, authHandler)

	return app, nil
}
