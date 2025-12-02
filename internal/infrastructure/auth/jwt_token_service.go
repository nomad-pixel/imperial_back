package auth

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type jwtTokenService struct {
	accessSecret         string
	refreshSecret        string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewJWTTokenService creates a new JWT token service with configuration
func NewJWTTokenService(accessSecret, refreshSecret string, accessDuration, refreshDuration time.Duration) ports.TokenService {
	return &jwtTokenService{
		accessSecret:         accessSecret,
		refreshSecret:        refreshSecret,
		accessTokenDuration:  accessDuration,
		refreshTokenDuration: refreshDuration,
	}
}

func (s *jwtTokenService) GenerateTokens(user *entities.User) (*entities.Tokens, error) {
	now := time.Now()

	// Generate access token with accessSecret
	accessClaims := jwt.MapClaims{
		"sub":   strconv.FormatInt(user.ID, 10),
		"email": user.Email,
		"typ":   "access",
		"exp":   now.Add(s.accessTokenDuration).Unix(),
		"iat":   now.Unix(),
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString([]byte(s.accessSecret))
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "Ошибка генерации access токена")
	}

	// Generate refresh token with refreshSecret
	refreshClaims := jwt.MapClaims{
		"sub": strconv.FormatInt(user.ID, 10),
		"typ": "refresh",
		"exp": now.Add(s.refreshTokenDuration).Unix(),
		"iat": now.Unix(),
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString([]byte(s.refreshSecret))
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "Ошибка генерации refresh токена")
	}

	tokens := &entities.Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	return tokens, nil
}

func (s *jwtTokenService) ValidateAccessToken(tokenStr string) (int64, error) {
	// Validate with accessSecret
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.New(apperrors.ErrCodeUnauthorized, "Invalid signing method")
		}
		return []byte(s.accessSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, apperrors.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, apperrors.ErrUnauthorized
	}

	if typ, _ := claims["typ"].(string); typ != "access" {
		return 0, apperrors.ErrUnauthorized
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, apperrors.ErrUnauthorized
	}
	id, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return 0, apperrors.ErrUnauthorized
	}
	return id, nil
}

func (s *jwtTokenService) RefreshAccessToken(refreshTokenStr string) (string, error) {
	// Validate refresh token with refreshSecret
	token, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.New(apperrors.ErrCodeUnauthorized, "Invalid signing method")
		}
		return []byte(s.refreshSecret), nil
	})
	if err != nil || !token.Valid {
		return "", apperrors.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", apperrors.ErrUnauthorized
	}
	if typ, _ := claims["typ"].(string); typ != "refresh" {
		return "", apperrors.ErrUnauthorized
	}
	sub, ok := claims["sub"].(string)
	if !ok {
		return "", apperrors.ErrUnauthorized
	}

	// Generate new access token with accessSecret
	now := time.Now()
	accessClaims := jwt.MapClaims{
		"sub": sub,
		"typ": "access",
		"exp": now.Add(s.accessTokenDuration).Unix(),
		"iat": now.Unix(),
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString([]byte(s.accessSecret))
	if err != nil {
		return "", apperrors.New(apperrors.ErrCodeInternal, "Ошибка генерации access токена")
	}
	return accessToken, nil
}
