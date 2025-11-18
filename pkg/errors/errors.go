package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorCode представляет тип ошибки
type ErrorCode string

const (
	// Клиентские ошибки (4xx)
	ErrCodeBadRequest   ErrorCode = "BAD_REQUEST"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"
	ErrCodeConflict     ErrorCode = "CONFLICT"
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"

	// Серверные ошибки (5xx)
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabase ErrorCode = "DATABASE_ERROR"
	ErrCodeExternal ErrorCode = "EXTERNAL_SERVICE_ERROR"
)

// AppError представляет ошибку приложения
type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
	Err        error                  `json:"-"`
}

// Error реализует интерфейс error
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap позволяет использовать errors.Is и errors.As
func (e *AppError) Unwrap() error {
	return e.Err
}

// New создает новую ошибку приложения
func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getStatusCode(code),
	}
}

// Wrap оборачивает существующую ошибку
func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getStatusCode(code),
		Err:        err,
	}
}

// WithDetails добавляет детали к ошибке
func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// getStatusCode возвращает HTTP статус код для типа ошибки
func getStatusCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest, ErrCodeValidation, ErrCodeInvalidInput:
		return http.StatusBadRequest
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeForbidden:
		return http.StatusForbidden
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	case ErrCodeDatabase, ErrCodeInternal, ErrCodeExternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// AsAppError пытается привести error к *AppError
func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// Предопределенные ошибки для частых случаев
var (
	ErrUserNotFound          = New(ErrCodeNotFound, "Пользователь не найден")
	ErrUserAlreadyExists     = New(ErrCodeConflict, "Пользователь уже существует")
	ErrInvalidCredentials    = New(ErrCodeUnauthorized, "Неверные учетные данные")
	ErrInvalidEmail          = New(ErrCodeValidation, "Неверный формат email")
	ErrPasswordTooShort      = New(ErrCodeValidation, "Пароль слишком короткий")
	ErrUnauthorized          = New(ErrCodeUnauthorized, "Требуется авторизация")
	ErrForbidden             = New(ErrCodeForbidden, "Доступ запрещен")
	ErrVerifyCodeNotFound    = New(ErrCodeNotFound, "Код верификации не найден")
	ErrVerifyCodeAlreadyUsed = New(ErrCodeConflict, "Код верификации уже использован")
	ErrVerifyCodeExpired     = New(ErrCodeValidation, "Код верификации истёк")
)
