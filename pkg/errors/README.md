# Обработка ошибок

Этот пакет предоставляет централизованную систему обработки ошибок для приложения.

## Архитектура

### AppError

Основная структура для представления ошибок приложения:

```go
type AppError struct {
    Code       ErrorCode              // Код ошибки
    Message    string                 // Сообщение для пользователя
    Details    map[string]interface{} // Дополнительные детали
    StatusCode int                    // HTTP статус код
    Err        error                  // Оригинальная ошибка
}
```

### ErrorCode

Типизированные коды ошибок:

**Клиентские ошибки (4xx):**
- `ErrCodeBadRequest` - Неверный запрос
- `ErrCodeUnauthorized` - Требуется авторизация
- `ErrCodeForbidden` - Доступ запрещен
- `ErrCodeNotFound` - Ресурс не найден
- `ErrCodeConflict` - Конфликт (например, дубликат)
- `ErrCodeValidation` - Ошибка валидации
- `ErrCodeInvalidInput` - Неверные входные данные

**Серверные ошибки (5xx):**
- `ErrCodeInternal` - Внутренняя ошибка
- `ErrCodeDatabase` - Ошибка базы данных
- `ErrCodeExternal` - Ошибка внешнего сервиса

## Использование

### В Usecases

```go
import "github.com/nomad-pixel/imperial/pkg/errors"

func (uc *MyUsecase) Execute(ctx context.Context) error {
    // Использование предопределенных ошибок
    if user == nil {
        return errors.ErrUserNotFound
    }

    // Создание новой ошибки
    if len(password) < 8 {
        return errors.New(errors.ErrCodeValidation, "Пароль слишком короткий")
    }

    // Оборачивание существующей ошибки
    if err != nil {
        return errors.Wrap(err, errors.ErrCodeDatabase, "Ошибка при сохранении")
    }

    // Добавление деталей
    return errors.New(errors.ErrCodeValidation, "Неверные данные").
        WithDetails("field", "email").
        WithDetails("value", email)
}
```

### В Handlers

```go
func (h *Handler) MyEndpoint(c *gin.Context) {
    result, err := h.usecase.Execute(c.Request.Context())
    if err != nil {
        // Просто передаем ошибку в middleware
        _ = c.Error(err)
        return
    }

    c.JSON(http.StatusOK, result)
}
```

### Middleware автоматически обработает ошибку

Middleware `ErrorHandler` автоматически:
1. Определяет тип ошибки (AppError или обычная)
2. Логирует серверные ошибки (5xx)
3. Возвращает правильный HTTP статус и JSON ответ

## Предопределенные ошибки

```go
var (
    ErrUserNotFound      = New(ErrCodeNotFound, "Пользователь не найден")
    ErrUserAlreadyExists = New(ErrCodeConflict, "Пользователь уже существует")
    ErrInvalidCredentials = New(ErrCodeUnauthorized, "Неверные учетные данные")
    ErrInvalidEmail      = New(ErrCodeValidation, "Неверный формат email")
    ErrPasswordTooShort  = New(ErrCodeValidation, "Пароль слишком короткий")
    ErrUnauthorized      = New(ErrCodeUnauthorized, "Требуется авторизация")
    ErrForbidden         = New(ErrCodeForbidden, "Доступ запрещен")
)
```

## Формат ответа

Все ошибки возвращаются в едином формате:

```json
{
  "code": "VALIDATION_ERROR",
  "message": "Неверный формат email",
  "details": {
    "field": "email",
    "value": "invalid-email"
  }
}
```

## Примеры

### Пример 1: Валидация

```go
func (uc *SignUpUsecase) validateEmail(email string) error {
    if !strings.Contains(email, "@") {
        return errors.ErrInvalidEmail.
            WithDetails("email", email)
    }
    return nil
}
```

### Пример 2: Обработка ошибок БД

```go
func (uc *Usecase) handleDBError(err error) error {
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        if pgErr.Code == "23505" {
            return errors.ErrUserAlreadyExists
        }
    }
    return errors.Wrap(err, errors.ErrCodeDatabase, "Ошибка БД")
}
```

### Пример 3: В Handler

```go
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        _ = c.Error(errors.Wrap(err, errors.ErrCodeValidation, "Неверный формат"))
        return
    }

    user, err := h.usecase.Execute(c.Request.Context(), req)
    if err != nil {
        _ = c.Error(err)
        return
    }

    c.JSON(http.StatusCreated, user)
}
```
