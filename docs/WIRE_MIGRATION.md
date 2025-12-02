# Migration to Wire DI

This document shows the before/after comparison of migrating from manual DI to Wire.

## 📊 Comparison

### Before (Manual DI)

**File size:** `internal/di/di.go` - **~250 lines**

```go
func InitializeApp(ctx context.Context, dbURL string) (*App, error) {
    // 1. Initialize database
    db, err := pgxpool.New(ctx, dbURL)
    if err != nil {
        return nil, err
    }

    // 2. Create all repositories (8+ lines)
    userRepo := repositories.NewUserRepositoryImpl(db)
    verifyCodeRepo := repositories.NewVerifyCodeRepositoryImpl(db)
    carCategoryRepo := repositories.NewCarCategoryRepositoryImpl(db)
    carTagRepo := repositories.NewCarTagRepositoryImpl(db)
    carMarkRepo := repositories.NewCarMarkRepositoryImpl(db)
    carRepo := repositories.NewCarRepositoryImpl(db)
    carImageRepo := repositories.NewCarImageRepositoryImpl(db)
    celebrityRepo := repositories.NewCelebrityRepositoryImpl(db)

    // 3. Initialize services (~30 lines)
    emailConfig := config.LoadEmailConfig()
    var emailService ports.EmailService
    if emailConfig.Provider == "smtp" {
        // ... SMTP setup
    } else {
        emailService = email.NewConsoleEmailService()
    }
    tokenSvc := token.NewJWTTokenService()
    imageService, err := imageSvc.NewFileImageService("./uploads", "http://localhost:8080/uploads")

    // 4. Create all use cases (~50 lines)
    signUpUsecase := authUsecase.NewSignUpUsecase(userRepo)
    sendEmailVerificationUsecase := authUsecase.NewSendEmailVerificationUsecase(userRepo, verifyCodeRepo, emailService)
    // ... 40+ more use cases

    // 5. Create all handlers (~30 lines)
    authHandler := auth.NewAuthHandler(signUpUsecase, sendEmailVerificationUsecase, ...)
    carHandler := car.NewCarHandler(createCarUsecase, deleteCarUsecase, ...)
    // ... 7+ more handlers

    // 6. Create app with 40+ parameters (~50 lines)
    app := NewApp(
        db,
        signUpUsecase,
        sendEmailVerificationUsecase,
        confirmEmailVerificationUsecase,
        signInUsecase,
        authHandler,
        tokenSvc,
        celebrityHandler,
        createCarUsecase,
        deleteCarUsecase,
        updateCarUsecase,
        // ... 35+ more parameters
    )

    return app, nil
}
```

**App struct:** ~70 fields

```go
type App struct {
    DB                 *pgxpool.Pool
    AuthHandler        *auth.AuthHandler
    CarHandler         *car.CarHandler
    // ... 7 more handlers

    // All use cases exposed (40+ fields)
    SignUpUsecase                   authUsecase.SignUpUsecase
    SendEmailVerificationUsecase    authUsecase.SendEmailVerificationUsecase
    CreateCarUsecase                carUsecase.CreateCarUsecase
    DeleteCarUsecase                carUsecase.DeleteCarUsecase
    // ... 35+ more use cases
}
```

### After (Wire DI)

**Total lines across files:** ~120 lines (52% reduction!)

**providers.go** (~60 lines):
```go
var ProviderSet = wire.NewSet(
    ProvideDatabase,
    ProvideEmailService,
    ProvideTokenService,
    ProvideImageService,

    // Repositories
    ProvideUserRepository,
    ProvideVerifyCodeRepository,
    // ... 6 more

    // Use case sets
    AuthUsecaseSet,
    CarUsecaseSet,
    CelebrityUsecaseSet,

    // Handlers
    HandlerSet,

    NewApp,
)

func ProvideDatabase(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
    return pgxpool.New(ctx, dbURL)
}

func ProvideEmailService() (ports.EmailService, error) {
    emailConfig := config.LoadEmailConfig()
    // ... implementation
}

// ... 8 more simple provider functions
```

**usecases.go** (~30 lines):
```go
var AuthUsecaseSet = wire.NewSet(
    authUsecase.NewSignUpUsecase,
    authUsecase.NewSendEmailVerificationUsecase,
    authUsecase.NewConfirmEmailVerificationUsecase,
    authUsecase.NewSignInUsecase,
    authUsecase.NewRefreshTokenUsecase,
)

var CarUsecaseSet = wire.NewSet(
    // Car CRUD
    carUsecase.NewCreateCarUsecase,
    carUsecase.NewDeleteCarUsecase,
    // ... all car use cases
)

var CelebrityUsecaseSet = wire.NewSet(
    celebrityUsecase.NewCreateCelebrityUsecase,
    // ... all celebrity use cases
)
```

**handlers.go** (~10 lines):
```go
var HandlerSet = wire.NewSet(
    auth.NewAuthHandler,
    car.NewCarHandler,
    carImage.NewCarImageHandler,
    carTag.NewCarTagHandler,
    carMark.NewCarMarkHandler,
    carCategory.NewCarCategoryHandler,
    celebrity.NewCelebrityHandler,
)
```

**wire.go** (~10 lines):
```go
//go:build wireinject
// +build wireinject

func InitializeApp(ctx context.Context, dbURL string) (*App, error) {
    wire.Build(ProviderSet)
    return nil, nil  // Wire generates the actual implementation
}
```

**app.go** - Simplified from 70+ fields to **10 fields**:
```go
type App struct {
    DB                 *pgxpool.Pool
    TokenService       ports.TokenService
    AuthHandler        *auth.AuthHandler
    CarHandler         *car.CarHandler
    CarImageHandler    *carImage.CarImageHandler
    CarTagHandler      *carTag.CarTagHandler
    CarMarkHandler     *carMark.CarMarkHandler
    CarCategoryHandler *carCategory.CarCategoryHandler
    CelebrityHandler   *celebrity.CelebrityHandler
}

func NewApp(
    db *pgxpool.Pool,
    tokenSvc ports.TokenService,
    authHandler *auth.AuthHandler,
    carHandler *car.CarHandler,
    carImageHandler *carImage.CarImageHandler,
    carTagHandler *carTag.CarTagHandler,
    carMarkHandler *carMark.CarMarkHandler,
    carCategoryHandler *carCategory.CarCategoryHandler,
    celebrityHandler *celebrity.CelebrityHandler,
) *App {
    return &App{
        DB:                 db,
        TokenService:       tokenSvc,
        AuthHandler:        authHandler,
        CarHandler:         carHandler,
        CarImageHandler:    carImageHandler,
        CarTagHandler:      carTagHandler,
        CarMarkHandler:     carMarkHandler,
        CarCategoryHandler: carCategoryHandler,
        CelebrityHandler:   celebrityHandler,
    }
}
```

**wire_gen.go** (Auto-generated - 90 lines):
```go
// Code generated by Wire. DO NOT EDIT.

func InitializeApp(ctx context.Context, dbURL string) (*App, error) {
    pool, err := ProvideDatabase(ctx, dbURL)
    if err != nil {
        return nil, err
    }
    tokenService := ProvideTokenService()
    userRepository := ProvideUserRepository(pool)
    signUpUsecase := usecases.NewSignUpUsecase(userRepository)
    // ... Wire automatically wires everything correctly
    app := NewApp(pool, tokenService, authHandler, carHandler, ...)
    return app, nil
}
```

## 📈 Benefits

### 1. **Code Reduction**
- **Before:** 250+ lines of manual wiring
- **After:** ~120 lines of declarations
- **Savings:** ~52% less code to maintain

### 2. **Better Organization**
- **Before:** One giant function with everything
- **After:** Organized by concern (providers, use cases, handlers)

### 3. **Type Safety**
- **Before:** Runtime errors if dependencies missing
- **After:** Compile-time validation of entire dependency graph

### 4. **Easier to Add Dependencies**

**Before - Adding a new feature:**
```go
// 1. Add repository (line 40)
articleRepo := repositories.NewArticleRepositoryImpl(db)

// 2. Add use cases (line 100)
createArticleUsecase := articleUsecase.NewCreateArticleUsecase(articleRepo)
getArticleUsecase := articleUsecase.NewGetArticleUsecase(articleRepo)

// 3. Add handler (line 150)
articleHandler := article.NewArticleHandler(createArticleUsecase, getArticleUsecase)

// 4. Add to App struct (app.go line 70)
ArticleHandler *article.ArticleHandler

// 5. Add to NewApp parameters (app.go line 85)
articleHandler *article.ArticleHandler,

// 6. Add to NewApp call (di.go line 200)
articleHandler,

// 7. Add to App initialization (app.go line 120)
ArticleHandler: articleHandler,
```

**After - Adding a new feature:**
```go
// 1. Add to providers.go (1 line)
ProvideArticleRepository,

// 2. Add to usecases.go (3 lines)
var ArticleUsecaseSet = wire.NewSet(
    article.NewCreateArticleUsecase,
    article.NewGetArticleUsecase,
)

// 3. Add to handlers.go (1 line)
article.NewArticleHandler,

// 4. Add to app.go (3 lines in struct, 3 lines in constructor)
ArticleHandler *article.ArticleHandler

// 5. Run: make wire
```

### 5. **Clear Dependency Graph**

Wire generates visual dependency graph:
```bash
cd internal/di && wire show
```

Shows exactly how dependencies flow:
```
App
├── *pgxpool.Pool (from ProvideDatabase)
├── TokenService (from ProvideTokenService)
├── *AuthHandler
│   ├── SignUpUsecase
│   │   └── UserRepository (from ProvideUserRepository)
│   ├── SignInUsecase
│   │   ├── UserRepository
│   │   └── TokenService
│   └── ...
└── ...
```

### 6. **No Runtime Overhead**

- **Before:** Function call overhead at startup
- **After:** Same generated code as manual, but verified at compile time

## 🚀 Adding a New Domain (Complete Example)

Let's say we want to add a "Comments" feature:

### Step 1: Create the domain code (same as before)

```go
// internal/domain/entities/comment.go
type Comment struct {
    ID        int64
    ArticleID int64
    UserID    int64
    Text      string
    CreatedAt time.Time
}

// internal/domain/ports/comment_repository.go
type CommentRepository interface {
    Create(ctx context.Context, comment *entities.Comment) error
    GetByArticleID(ctx context.Context, articleID int64) ([]*entities.Comment, error)
}

// internal/domain/usecases/comment/create_comment.go
func NewCreateCommentUsecase(repo ports.CommentRepository) CreateCommentUsecase { ... }

// internal/domain/usecases/comment/get_comments.go
func NewGetCommentsUsecase(repo ports.CommentRepository) GetCommentsUsecase { ... }

// internal/infrastructure/repositories/comment_repository_impl.go
func NewCommentRepositoryImpl(db *pgxpool.Pool) ports.CommentRepository { ... }

// internal/interfaces/http/comment/handler.go
func NewCommentHandler(
    createUsecase usecases.CreateCommentUsecase,
    getUsecase usecases.GetCommentsUsecase,
) *CommentHandler { ... }
```

### Step 2: Wire it up (new simplified way)

**providers.go:**
```go
func ProvideCommentRepository(db *pgxpool.Pool) ports.CommentRepository {
    return repositories.NewCommentRepositoryImpl(db)
}

var ProviderSet = wire.NewSet(
    // ... existing
    ProvideCommentRepository,  // Add this
    CommentUsecaseSet,         // Add this
)
```

**usecases.go:**
```go
var CommentUsecaseSet = wire.NewSet(
    comment.NewCreateCommentUsecase,
    comment.NewGetCommentsUsecase,
)
```

**handlers.go:**
```go
var HandlerSet = wire.NewSet(
    // ... existing
    comment.NewCommentHandler,  // Add this
)
```

**app.go:**
```go
type App struct {
    // ... existing
    CommentHandler *comment.CommentHandler  // Add this
}

func NewApp(
    // ... existing
    commentHandler *comment.CommentHandler,  // Add this
) *App {
    return &App{
        // ... existing
        CommentHandler: commentHandler,  // Add this
    }
}
```

### Step 3: Generate and done!

```bash
make wire
```

Wire generates all the wiring code automatically with full type checking!

## 🎯 Migration Checklist

- [x] Install Wire (`go get github.com/google/wire/cmd/wire`)
- [x] Create `providers.go` with provider functions
- [x] Create `usecases.go` with use case sets
- [x] Create `handlers.go` with handler set
- [x] Simplify `app.go` (remove use case fields)
- [x] Create `wire.go` with injector declaration
- [x] Generate `wire_gen.go` (`make wire`)
- [x] Remove old `di.go`
- [x] Update Makefile with `wire` command
- [x] Test build (`go build ./cmd/api`)
- [x] Document in README

## 📚 Next Steps

1. **Add validation** - Use Wire to inject validators
2. **Add caching** - Wire up cache layers
3. **Add monitoring** - Inject metrics collectors
4. **Environment configs** - Provider functions for different envs

## 🔍 Troubleshooting

### Build errors after migration?

Make sure old `di.go` is deleted:
```bash
rm internal/di/di.go.old
```

### Wire generation fails?

Check that all constructors return the right types:
```bash
cd internal/di && wire check
```

### Missing dependency error?

Add the provider to `ProviderSet` in `providers.go`

## 🎉 Result

**Before:**
- 250+ lines of repetitive wiring code
- Easy to make mistakes
- Hard to maintain
- Runtime errors possible

**After:**
- 120 lines of declarative providers
- Compile-time type checking
- Easy to maintain and extend
- Impossible to wire incorrectly

**Wire saves time and prevents bugs!** 🚀
