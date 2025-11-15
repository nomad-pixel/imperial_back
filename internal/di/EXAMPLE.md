# Пример добавления новых зависимостей

## Сценарий: Добавление функционала для работы с продуктами

### Шаг 1: Создайте Entity

```go
// internal/domain/entities/product.go
package entities

type Product struct {
    ID        string
    Name      string
    Price     float64
    CreatedAt time.Time
}
```

### Шаг 2: Создайте Port (интерфейс репозитория)

```go
// internal/domain/ports/product_repository.go
package ports

type ProductRepository interface {
    CreateProduct(ctx context.Context, name string, price float64) (*entities.Product, error)
    GetProductByID(ctx context.Context, id string) (*entities.Product, error)
}
```

### Шаг 3: Реализуйте репозиторий

```go
// internal/infrastructure/repositories/product_repository_impl.go
package repositories

type ProductRepositoryImpl struct {
    db *pgxpool.Pool
}

func NewProductRepositoryImpl(db *pgxpool.Pool) ports.ProductRepository {
    return &ProductRepositoryImpl{db: db}
}

func (r *ProductRepositoryImpl) CreateProduct(ctx context.Context, name string, price float64) (*entities.Product, error) {
    // implementation
}
```

### Шаг 4: Создайте UseCase

```go
// internal/domain/usecases/create_product_usecase.go
package usecases

type CreateProductUsecase struct {
    productRepository ports.ProductRepository
}

func NewCreateProductUsecase(productRepository ports.ProductRepository) *CreateProductUsecase {
    return &CreateProductUsecase{productRepository: productRepository}
}

func (uc *CreateProductUsecase) Execute(ctx context.Context, name string, price float64) (*entities.Product, error) {
    return uc.productRepository.CreateProduct(ctx, name, price)
}
```

### Шаг 5: Создайте Handler

```go
// internal/interfaces/http/product/handler.go
package product

type ProductHandler struct {
    createProductUsecase *usecases.CreateProductUsecase
}

func NewProductHandler(createProductUsecase *usecases.CreateProductUsecase) *ProductHandler {
    return &ProductHandler{createProductUsecase: createProductUsecase}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
    // implementation
}
```

### Шаг 6: Обновите DI контейнер

```go
// internal/di/wire.go
func InitializeApp(ctx context.Context, dbURL string) (*App, error) {
    db, err := pgxpool.New(ctx, dbURL)
    if err != nil {
        return nil, err
    }

    // Repositories
    userRepo := repositories.NewUserRepositoryImpl(db)
    productRepo := repositories.NewProductRepositoryImpl(db)  // +++ ДОБАВЛЕНО

    // Usecases
    signUpUsecase := usecases.NewSignUpUsecase(userRepo)
    createProductUsecase := usecases.NewCreateProductUsecase(productRepo)  // +++ ДОБАВЛЕНО

    // Handlers
    authHandler := auth.NewAuthHandler(signUpUsecase)
    productHandler := product.NewProductHandler(createProductUsecase)  // +++ ДОБАВЛЕНО

    // App
    app := NewApp(db, authHandler, productHandler)  // +++ ОБНОВЛЕНО

    return app, nil
}
```

### Шаг 7: Обновите структуру App

```go
// internal/di/app.go
type App struct {
    DB             *pgxpool.Pool
    AuthHandler    *auth.AuthHandler
    ProductHandler *product.ProductHandler  // +++ ДОБАВЛЕНО
}

func NewApp(
    db *pgxpool.Pool,
    authHandler *auth.AuthHandler,
    productHandler *product.ProductHandler,  // +++ ДОБАВЛЕНО
) *App {
    return &App{
        DB:             db,
        AuthHandler:    authHandler,
        ProductHandler: productHandler,  // +++ ДОБАВЛЕНО
    }
}
```

### Шаг 8: Зарегистрируйте роуты в main.go

```go
// cmd/api/main.go
func main() {
    // ... initialization code ...

    server := gin.Default()

    auth.RegisterRoutes(server, app.AuthHandler)
    product.RegisterRoutes(server, app.ProductHandler)  // +++ ДОБАВЛЕНО

    server.Run(":8080")
}
```

## Преимущества такого подхода

1. **Централизованная инициализация** - все зависимости создаются в одном месте
2. **Явные зависимости** - легко увидеть, что от чего зависит
3. **Тестируемость** - легко подменить зависимости при тестировании
4. **Типобезопасность** - компилятор проверяет правильность связей
5. **Простота понимания** - четкий порядок инициализации без "магии"
