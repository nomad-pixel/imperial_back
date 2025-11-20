package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/car"
)

type App struct {
	DB                              *pgxpool.Pool
	AuthHandler                     *auth.AuthHandler
	CarHandler                      *car.CarHandler
	SignUpUsecase                   usecasePorts.SignUpUsecase
	SignInUsecase                   usecasePorts.SignInUsecase
	SendEmailVerificationUsecase    usecasePorts.SendEmailVerificationUsecase
	ConfirmEmailVerificationUsecase usecasePorts.ConfirmEmailVerificationUsecase
	TokenService                    ports.TokenService

	CreateCarUsecase usecasePorts.CreateCarUsecase
}

func NewApp(
	db *pgxpool.Pool,
	signUpUsecase usecasePorts.SignUpUsecase,
	sendEmailVerificationUsecase usecasePorts.SendEmailVerificationUsecase,
	confirmEmailVerificationUsecase usecasePorts.ConfirmEmailVerificationUsecase,
	signInUsecase usecasePorts.SignInUsecase,
	authHandler *auth.AuthHandler,
	tokenSvc ports.TokenService,

	createCarUsecase usecasePorts.CreateCarUsecase,
	carHandler *car.CarHandler,

) *App {
	return &App{
		DB:                              db,
		AuthHandler:                     authHandler,
		SignUpUsecase:                   signUpUsecase,
		ConfirmEmailVerificationUsecase: confirmEmailVerificationUsecase,
		SendEmailVerificationUsecase:    sendEmailVerificationUsecase,
		SignInUsecase:                   signInUsecase,
		TokenService:                    tokenSvc,

		CarHandler:       carHandler,
		CreateCarUsecase: createCarUsecase,
	}
}

func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
