package di

import (
	"github.com/jackc/pgx/v5/pgxpool"
	usecasePorts "github.com/nomad-pixel/imperial/internal/domain/usecases"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
)

// App содержит все зависимости приложения
type App struct {
	DB            *pgxpool.Pool
	AuthHandler   *auth.AuthHandler
	SignUpUsecase usecasePorts.SignUpUsecase
}

// NewApp создает новый экземпляр приложения
func NewApp(
	db *pgxpool.Pool,
	signUpUsecase usecasePorts.SignUpUsecase,
	authHandler *auth.AuthHandler,
) *App {
	return &App{
		DB:            db,
		AuthHandler:   authHandler,
		SignUpUsecase: signUpUsecase,
	}
}

// Close закрывает все ресурсы приложения
func (a *App) Close() {
	if a.DB != nil {
		a.DB.Close()
	}
}
