package di

import (
	"context"
	"fmt"
	"log"

	"github.com/google/wire"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nomad-pixel/imperial/internal/config"
	"github.com/nomad-pixel/imperial/internal/domain/ports"
	token "github.com/nomad-pixel/imperial/internal/infrastructure/auth"
	"github.com/nomad-pixel/imperial/internal/infrastructure/email"
	imageSvc "github.com/nomad-pixel/imperial/internal/infrastructure/image"
	postgres "github.com/nomad-pixel/imperial/internal/infrastructure/postgres"
)

var ProviderSet = wire.NewSet(
	// Config provider
	ProvideConfig,

	// Infrastructure providers
	ProvideDatabase,
	ProvideEmailService,
	ProvideTokenService,
	ProvideImageService,

	// Repository providers
	ProvideUserRepository,
	ProvideVerifyCodeRepository,
	ProvideCarRepository,
	ProvideCarCategoryRepository,
	ProvideCarTagRepository,
	ProvideCarMarkRepository,
	ProvideCarImageRepository,
	ProvideCelebrityRepository,
	ProvideLeadRepository,
	ProvideDriverRepository,

	// Use case providers (imported from other files)
	AuthUsecaseSet,
	CarUsecaseSet,
	CelebrityUsecaseSet,
	LeadUsecaseSet,
	DriverUsecaseSet,

	// Handler providers
	HandlerSet,

	// App provider
	NewApp,
)

// ProvideConfig loads and validates application configuration
func ProvideConfig() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

func ProvideDatabase(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}
	poolConfig.MaxConns = cfg.Database.MaxConns
	poolConfig.MinConns = cfg.Database.MinConns
	poolConfig.MaxConnLifetime = cfg.Database.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.Database.MaxConnIdleTime

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test connection
	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func ProvideEmailService(cfg *config.Config) (ports.EmailService, error) {
	if cfg.Email.Provider == "smtp" {
		log.Printf("Initializing SMTP email service (host: %s, port: %d)", cfg.Email.SMTP.Host, cfg.Email.SMTP.Port)

		smtpService, err := email.NewSMTPEmailService(email.SMTPConfig{
			Host:     cfg.Email.SMTP.Host,
			Port:     fmt.Sprintf("%d", cfg.Email.SMTP.Port),
			Username: cfg.Email.SMTP.Username,
			Password: cfg.Email.SMTP.Password,
			From:     cfg.Email.SMTP.From,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to initialize SMTP service: %w", err)
		}
		return smtpService, nil
	}

	log.Println("Using console email service (emails will be printed to console)")
	return email.NewConsoleEmailService(), nil
}

func ProvideTokenService(cfg *config.Config) ports.TokenService {
	log.Printf("Initializing JWT token service (access: %v, refresh: %v)",
		cfg.JWT.AccessTokenDuration, cfg.JWT.RefreshTokenDuration)
	return token.NewJWTTokenService(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTokenDuration,
		cfg.JWT.RefreshTokenDuration,
	)
}

func ProvideImageService(cfg *config.Config) (ports.ImageService, error) {
	log.Printf("Initializing file storage (path: %s, base URL: %s)", cfg.Storage.LocalPath, cfg.Storage.BaseURL)

	imageService, err := imageSvc.NewFileImageService(cfg.Storage.LocalPath, cfg.Storage.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize image service: %w", err)
	}
	return imageService, nil
}

func ProvideUserRepository(db *pgxpool.Pool) ports.UserRepository {
	return postgres.NewUserRepositoryImpl(db)
}

func ProvideVerifyCodeRepository(db *pgxpool.Pool) ports.VerifyCodeRepository {
	return postgres.NewVerifyCodeRepositoryImpl(db)
}

func ProvideCarRepository(db *pgxpool.Pool) ports.CarRepository {
	return postgres.NewCarRepositoryImpl(db)
}

func ProvideCarCategoryRepository(db *pgxpool.Pool) ports.CarCategoryRepository {
	return postgres.NewCarCategoryRepositoryImpl(db)
}

func ProvideCarTagRepository(db *pgxpool.Pool) ports.CarTagRepository {
	return postgres.NewCarTagRepositoryImpl(db)
}

func ProvideCarMarkRepository(db *pgxpool.Pool) ports.CarMarkRepository {
	return postgres.NewCarMarkRepositoryImpl(db)
}

func ProvideCarImageRepository(db *pgxpool.Pool) ports.CarImageRepository {
	return postgres.NewCarImageRepositoryImpl(db)
}

func ProvideCelebrityRepository(db *pgxpool.Pool) ports.CelebrityRepository {
	return postgres.NewCelebrityRepositoryImpl(db)
}

func ProvideLeadRepository(db *pgxpool.Pool) ports.LeadRepository {
	return postgres.NewLeadRepository(db)
}

func ProvideDriverRepository(db *pgxpool.Pool) ports.DriverRepository {
	return postgres.NewDriverRepository(db)
}
