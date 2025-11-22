package auth

import (
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type jwtTokenService struct {
}

func NewJWTTokenService() ports.TokenService {
	return &jwtTokenService{}
}

func (s *jwtTokenService) GenerateTokens(user *entities.User) (*entities.Tokens, error) {
	accessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if accessSecret == "" || refreshSecret == "" {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "Token secrets not configured in environment")
	}

	accessExpMinutes := 15
	if v := os.Getenv("ACCESS_TOKEN_EXPIRES_MINUTES"); v != "" {
		if iv, e := strconv.Atoi(v); e == nil && iv > 0 {
			accessExpMinutes = iv
		}
	}
	refreshExpDays := 7
	if v := os.Getenv("REFRESH_TOKEN_EXPIRES_DAYS"); v != "" {
		if iv, e := strconv.Atoi(v); e == nil && iv > 0 {
			refreshExpDays = iv
		}
	}

	accessClaims := jwt.MapClaims{
		"sub":   strconv.FormatInt(user.ID, 10),
		"email": user.Email,
		"typ":   "access",
		"exp":   time.Now().Add(time.Duration(accessExpMinutes) * time.Minute).Unix(),
		"iat":   time.Now().Unix(),
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString([]byte(accessSecret))
	if err != nil {
		return nil, apperrors.New(apperrors.ErrCodeInternal, "Ошибка генерации access токена")
	}

	refreshClaims := jwt.MapClaims{
		"sub": strconv.FormatInt(user.ID, 10),
		"typ": "refresh",
		"exp": time.Now().Add(time.Duration(refreshExpDays) * 24 * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err := refreshTokenObj.SignedString([]byte(refreshSecret))
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
	accessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if accessSecret == "" {
		return 0, apperrors.New(apperrors.ErrCodeInternal, "ACCESS_TOKEN_SECRET not configured")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.New(apperrors.ErrCodeUnauthorized, "Invalid signing method")
		}
		return []byte(accessSecret), nil
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
	refreshSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	accessSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if refreshSecret == "" || accessSecret == "" {
		return "", apperrors.New(apperrors.ErrCodeInternal, "Token secrets not configured in environment")
	}

	token, err := jwt.Parse(refreshTokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperrors.New(apperrors.ErrCodeUnauthorized, "Invalid signing method")
		}
		return []byte(refreshSecret), nil
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

	accessExpMinutes := 15
	if v := os.Getenv("ACCESS_TOKEN_EXPIRES_MINUTES"); v != "" {
		if iv, e := strconv.Atoi(v); e == nil && iv > 0 {
			accessExpMinutes = iv
		}
	}

	accessClaims := jwt.MapClaims{
		"sub": sub,
		"typ": "access",
		"exp": time.Now().Add(time.Duration(accessExpMinutes) * time.Minute).Unix(),
		"iat": time.Now().Unix(),
	}
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessTokenObj.SignedString([]byte(accessSecret))
	if err != nil {
		return "", apperrors.New(apperrors.ErrCodeInternal, "Ошибка генерации access токена")
	}
	return accessToken, nil
}
