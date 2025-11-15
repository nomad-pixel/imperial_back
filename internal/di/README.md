# Dependency Injection

Этот пакет содержит логику инициализации всех зависимостей приложения.

## Структура

### app.go
Содержит структуру `App`, которая объединяет все зависимости приложения:
- `DB` - подключение к базе данных
- `AuthHandler` - обработчик HTTP запросов для авторизации

### wire.go
Содержит функцию `InitializeApp`, которая инициализирует все зависимости в правильном порядке:
1. Database Pool
2. Repositories (UserRepository)
3. Usecases (SignUpUsecase)
4. Handlers (AuthHandler)

## Как добавить новую зависимость

1. Создайте новый компонент (repository, usecase, handler)
2. Добавьте его в функцию `InitializeApp` в правильном порядке
3. Если нужно, добавьте в структуру `App` для доступа из main.go

## Пример

```go
// Repositories
userRepo := repositories.NewUserRepositoryImpl(db)
productRepo := repositories.NewProductRepositoryImpl(db)

// Usecases
signUpUsecase := usecases.NewSignUpUsecase(userRepo)
createProductUsecase := usecases.NewCreateProductUsecase(productRepo)

// Handlers
authHandler := auth.NewAuthHandler(signUpUsecase)
productHandler := product.NewProductHandler(createProductUsecase)
```
