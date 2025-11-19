package ports

import "github.com/nomad-pixel/imperial/internal/domain/entities"

type TokenService interface {
	GenerateTokens(user *entities.User) (*entities.Tokens, error)
}
