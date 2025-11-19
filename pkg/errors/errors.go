package errors

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	ErrCodeBadRequest   ErrorCode = "BAD_REQUEST"            // HTTP 400
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"           // HTTP 401
	ErrCodeForbidden    ErrorCode = "FORBIDDEN"              // HTTP 403
	ErrCodeNotFound     ErrorCode = "NOT_FOUND"              // HTTP 404
	ErrCodeConflict     ErrorCode = "CONFLICT"               // HTTP 409
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"       // HTTP 400
	ErrCodeInvalidInput ErrorCode = "INVALID_INPUT"          // HTTP 400
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"         // HTTP 500
	ErrCodeDatabase     ErrorCode = "DATABASE_ERROR"         // HTTP 500
	ErrCodeExternal     ErrorCode = "EXTERNAL_SERVICE_ERROR" // HTTP 500

)

type AppError struct {
	Code       ErrorCode              `json:"code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	StatusCode int                    `json:"-"`
	Err        error                  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getStatusCode(code),
	}
}

func Wrap(err error, code ErrorCode, message string) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		StatusCode: getStatusCode(code),
		Err:        err,
	}
}

func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

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

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

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
	ErrUserNotVerified       = New(ErrCodeUnauthorized, "Пользователь не верифицирован")
)
