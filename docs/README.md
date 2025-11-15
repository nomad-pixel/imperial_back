# API Documentation

Эта директория содержит автоматически сгенерированную Swagger документацию для Imperial API.

## Генерация документации

Документация генерируется автоматически из комментариев в коде с использованием [swaggo/swag](https://github.com/swaggo/swag).

### Генерация

```bash
make swagger
```

Или напрямую:

```bash
~/go/bin/swag init -g cmd/api/main.go -o docs
```

## Просмотр документации

После запуска приложения, Swagger UI доступен по адресу:

```
http://localhost:8080/swagger/index.html
```

## Файлы

- `docs.go` - Go код с метаданными API
- `swagger.json` - JSON спецификация OpenAPI
- `swagger.yaml` - YAML спецификация OpenAPI

## Как документировать новые endpoints

### Пример аннотаций

```go
// SignUp godoc
// @Summary      Регистрация нового пользователя
// @Description  Создает нового пользователя с указанным email и паролем
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body SignUpRequest true "Данные для регистрации"
// @Success      201 {object} SignUpResponse "Пользователь успешно зарегистрирован"
// @Failure      400 {object} ErrorResponse "Неверные данные запроса"
// @Failure      500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router       /api/v1/auth/sign-up [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
    // implementation
}
```

### Структуры DTO

Используйте теги для примеров и валидации:

```go
type SignUpRequest struct {
    Email    string `json:"email" binding:"required,email" example:"user@example.com"`
    Password string `json:"password" binding:"required,min=8" example:"password123"`
}
```

## После изменений

1. Добавьте аннотации к новым endpoints
2. Запустите `make swagger` для регенерации документации
3. Пересоберите приложение `make build`
4. Проверьте документацию в браузере
