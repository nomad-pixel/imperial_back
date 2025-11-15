# Система обработки ошибок

## Обзор

В проекте реализована централизованная система обработки ошибок с использованием:
- Типизированных кодов ошибок
- Middleware для автоматической обработки
- Предопределенных ошибок для частых случаев
- Структурированного формата ответов

## Компоненты

### 1. Пакет `pkg/errors`

Содержит типы и функции для работы с ошибками:

```go
type AppError struct {
    Code       ErrorCode              // Тип ошибки
    Message    string                 // Сообщение пользователю
    Details    map[string]interface{} // Дополнительные данные
    StatusCode int                    // HTTP статус код
    Err        error                  // Оригинальная ошибка
}
```

### 2. Middleware `ErrorHandler`

Автоматически обрабатывает все ошибки из handlers:
- Распознает `AppError` и обычные errors
- Логирует серверные ошибки (5xx)
- Возвращает структурированный JSON ответ

### 3. Middleware `Recovery`

Ловит панику и преобразует в HTTP 500 ответ.

## Коды ошибок

### Клиентские (4xx)
- `BAD_REQUEST` - 400
- `UNAUTHORIZED` - 401
- `FORBIDDEN` - 403
- `NOT_FOUND` - 404
- `CONFLICT` - 409
- `VALIDATION_ERROR` - 400
- `INVALID_INPUT` - 400

### Серверные (5xx)
- `INTERNAL_ERROR` - 500
- `DATABASE_ERROR` - 500
- `EXTERNAL_SERVICE_ERROR` - 500

## Использование

### В Usecases

```go
func (uc *SignUpUsecase) Execute(ctx context.Context, email, password string) (*entities.User, error) {
    // Валидация
    if len(password) < 8 {
        return nil, errors.ErrPasswordTooShort
    }

    // Обработка ошибок БД
    user, err := uc.userRepository.CreateUser(ctx, email, password)
    if err != nil {
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgErr.Code == "23505" {
            return nil, errors.ErrUserAlreadyExists
        }
        return nil, errors.Wrap(err, errors.ErrCodeDatabase, "Ошибка создания пользователя")
    }

    return user, nil
}
```

### В Handlers

```go
func (h *AuthHandler) SignUp(c *gin.Context) {
    var req SignUpRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        _ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат"))
        return
    }

    user, err := h.signUpUsecase.Execute(c.Request.Context(), req.Email, req.Password)
    if err != nil {
        _ = c.Error(err)  // Middleware обработает автоматически
        return
    }

    c.JSON(http.StatusCreated, user)
}
```

## Формат ответа

Все ошибки возвращаются в едином JSON формате:

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Пароль слишком короткий",
  "details": {
    "min_length": 8,
    "actual_length": 5
  }
}
```

## Примеры ответов

### 400 Bad Request - Валидация
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Неверный формат email"
}
```

### 409 Conflict - Дубликат
```json
{
  "code": "CONFLICT",
  "message": "Пользователь уже существует"
}
```

### 500 Internal Server Error
```json
{
  "code": "DATABASE_ERROR",
  "message": "Ошибка при создании пользователя"
}
```

## Предопределенные ошибки

```go
errors.ErrUserNotFound       // 404
errors.ErrUserAlreadyExists  // 409
errors.ErrInvalidCredentials // 401
errors.ErrInvalidEmail       // 400
errors.ErrPasswordTooShort   // 400
errors.ErrUnauthorized       // 401
errors.ErrForbidden          // 403
```

## Best Practices

1. **В Usecases**: Используйте `AppError` для бизнес-логики
2. **В Handlers**: Просто передавайте ошибки в `c.Error()`
3. **В Repository**: Возвращайте стандартные errors, обрабатывайте в usecase
4. **Логирование**: Серверные ошибки логируются автоматически
5. **Детали**: Добавляйте детали через `.WithDetails()` для отладки

## Добавление новых ошибок

```go
// В pkg/errors/errors.go
const (
    ErrCodeCustom ErrorCode = "CUSTOM_ERROR"
)

var ErrCustom = New(ErrCodeCustom, "Описание ошибки")
```

## Тестирование

```go
func TestErrorHandling(t *testing.T) {
    err := errors.ErrUserNotFound

    // Проверка типа
    appErr, ok := errors.AsAppError(err)
    assert.True(t, ok)
    assert.Equal(t, errors.ErrCodeNotFound, appErr.Code)
    assert.Equal(t, http.StatusNotFound, appErr.StatusCode)
}
```
