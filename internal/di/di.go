package di

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/domain/usecases"
	"github.com/nomad-pixel/imperial/internal/infrastructure/repositories"
	"github.com/nomad-pixel/imperial/internal/interfaces/http/auth"
)

// InitializeApp инициализирует все зависимости приложения
func InitializeApp(ctx context.Context, dbURL string) (*App, error) {
	// Database
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	// Repositories
	userRepo := repositories.NewUserRepositoryImpl(db)

	// Usecases
	signUpUsecase := usecases.NewSignUpUsecase(userRepo)

	// Handlers
	authHandler := auth.NewAuthHandler(signUpUsecase)

	// App
	app := NewApp(db, signUpUsecase, authHandler)

	return app, nil
}
