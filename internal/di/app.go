package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
)

type App struct {
	DB                              *pgxpool.Pool
	AuthHandler                     *auth.AuthHandler
	SignUpUsecase                   usecasePorts.SignUpUsecase
	SignInUsecase                   usecasePorts.SignInUsecase
	SendEmailVerificationUsecase    usecasePorts.SendEmailVerificationUsecase
	ConfirmEmailVerificationUsecase usecasePorts.ConfirmEmailVerificationUsecase
}

func NewApp(
	db *pgxpool.Pool,
	signUpUsecase usecasePorts.SignUpUsecase,
	sendEmailVerificationUsecase usecasePorts.SendEmailVerificationUsecase,
	confirmEmailVerificationUsecase usecasePorts.ConfirmEmailVerificationUsecase,
	signInUsecase usecasePorts.SignInUsecase,
	authHandler *auth.AuthHandler,
) *App {
	return &App{
		DB:                              db,
		AuthHandler:                     authHandler,
		SignUpUsecase:                   signUpUsecase,
		ConfirmEmailVerificationUsecase: confirmEmailVerificationUsecase,
		SendEmailVerificationUsecase:    sendEmailVerificationUsecase,
		SignInUsecase:                   signInUsecase,
	}
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
