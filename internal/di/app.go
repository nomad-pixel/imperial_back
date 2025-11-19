package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
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
	TokenService                    ports.TokenService
}

func NewApp(
	db *pgxpool.Pool,
	signUpUsecase usecasePorts.SignUpUsecase,
	sendEmailVerificationUsecase usecasePorts.SendEmailVerificationUsecase,
	confirmEmailVerificationUsecase usecasePorts.ConfirmEmailVerificationUsecase,
	signInUsecase usecasePorts.SignInUsecase,
	authHandler *auth.AuthHandler,
	tokenSvc ports.TokenService,
) *App {
	return &App{
		DB:                              db,
		AuthHandler:                     authHandler,
		SignUpUsecase:                   signUpUsecase,
		ConfirmEmailVerificationUsecase: confirmEmailVerificationUsecase,
		SendEmailVerificationUsecase:    sendEmailVerificationUsecase,
		SignInUsecase:                   signInUsecase,
		TokenService:                    tokenSvc,
	}
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
