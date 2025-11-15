package usecases

import (
	"context"
	"strings"

	"github.com/nomad-pixel/imperial/internal/domain/entities"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	"github.com/nomad-pixel/imperial/pkg/errors"
)

type SignUpUsecase struct {
	userRepository ports.UserRepository
}

func NewSignUpUsecase(userRepository ports.UserRepository) *SignUpUsecase {
	return &SignUpUsecase{userRepository: userRepository}
}

func (uc *SignUpUsecase) Execute(ctx context.Context, email, password string) (*entities.User, error) {
	// Валидация входных данных
	if err := uc.validateInput(email, password); err != nil {
		return nil, err
	}

	// Проверка существования пользователя
	_, err := uc.userRepository.GetUserByEmail(ctx, email)
	if err == nil {
		// Пользователь найден - значит уже существует
		return nil, errors.ErrUserAlreadyExists
	}
	// Если ошибка не "пользователь не найден", значит проблема с БД
	if err != errors.ErrUserNotFound {
		return nil, err
	}

	// Создание пользователя
	user, err := uc.userRepository.CreateUser(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *SignUpUsecase) validateInput(email, password string) error {
	// Проверка email
	if email == "" || !strings.Contains(email, "@") {
		return errors.ErrInvalidEmail
	}

	// Проверка пароля
	if len(password) < 8 {
		return errors.ErrPasswordTooShort
	}

	return nil
}
