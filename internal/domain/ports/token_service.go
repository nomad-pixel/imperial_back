package ports

import "github.com/nomad-pixel/imperial/internal/domain/entities"

type TokenService interface {
	GenerateTokens(user *entities.User) (*entities.Tokens, error)
	ValidateAccessToken(token string) (int64, error)
	RefreshAccessToken(refreshToken string) (string, error)
}
