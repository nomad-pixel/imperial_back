package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Email    EmailConfig
	Storage  StorageConfig
	Server   ServerConfig
}

// AppConfig contains general application settings
type AppConfig struct {
	Name        string `envconfig:"APP_NAME" default:"Imperial"`
	Environment string `envconfig:"APP_ENV" default:"development"`
	Debug       bool   `envconfig:"APP_DEBUG" default:"false"`
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	URL             string        `envconfig:"DATABASE_URL" required:"true"`
	MaxConns        int32         `envconfig:"DB_MAX_CONNS" default:"25"`
	MinConns        int32         `envconfig:"DB_MIN_CONNS" default:"5"`
	MaxConnLifetime time.Duration `envconfig:"DB_MAX_CONN_LIFETIME" default:"1h"`
	MaxConnIdleTime time.Duration `envconfig:"DB_MAX_CONN_IDLE_TIME" default:"30m"`
}

// JWTConfig contains JWT token settings
type JWTConfig struct {
	AccessSecret         string        `envconfig:"JWT_ACCESS_SECRET" required:"true"`
	RefreshSecret        string        `envconfig:"JWT_REFRESH_SECRET" required:"true"`
	AccessTokenDuration  time.Duration `envconfig:"JWT_ACCESS_DURATION" default:"15m"`
	RefreshTokenDuration time.Duration `envconfig:"JWT_REFRESH_DURATION" default:"168h"` // 7 days
}

// EmailConfig contains email service settings
type EmailConfig struct {
	Provider string     `envconfig:"EMAIL_PROVIDER" default:"console"`
	SMTP     SMTPConfig `envconfig:"SMTP"`
}

// SMTPConfig contains SMTP server settings
type SMTPConfig struct {
	Host     string `envconfig:"SMTP_HOST" default:"smtp.gmail.com"`
	Port     int    `envconfig:"SMTP_PORT" default:"587"`
	Username string `envconfig:"SMTP_USER"`
	Password string `envconfig:"SMTP_PASS"`
	From     string `envconfig:"SMTP_FROM" default:"noreply@imperial.com"`
}

// StorageConfig contains file storage settings
type StorageConfig struct {
	Type      string `envconfig:"STORAGE_TYPE" default:"local"`
	LocalPath string `envconfig:"STORAGE_LOCAL_PATH" default:"./uploads"`
	BaseURL   string `envconfig:"STORAGE_BASE_URL" default:"http://localhost:8080/uploads"`
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Port            int           `envconfig:"SERVER_PORT" default:"8080"`
	ReadTimeout     time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"10s"`
	WriteTimeout    time.Duration `envconfig:"SERVER_WRITE_TIMEOUT" default:"10s"`
	ShutdownTimeout time.Duration `envconfig:"SERVER_SHUTDOWN_TIMEOUT" default:"5s"`
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	var cfg Config

	// Load App config
	if err := envconfig.Process("", &cfg.App); err != nil {
		return nil, fmt.Errorf("failed to load app config: %w", err)
	}

	// Load Database config
	if err := envconfig.Process("", &cfg.Database); err != nil {
		return nil, fmt.Errorf("failed to load database config: %w", err)
	}

	// Load JWT config
	if err := envconfig.Process("", &cfg.JWT); err != nil {
		return nil, fmt.Errorf("failed to load JWT config: %w", err)
	}

	// Load Email config
	if err := envconfig.Process("", &cfg.Email); err != nil {
		return nil, fmt.Errorf("failed to load email config: %w", err)
	}

	// Load Storage config
	if err := envconfig.Process("", &cfg.Storage); err != nil {
		return nil, fmt.Errorf("failed to load storage config: %w", err)
	}

	// Load Server config
	if err := envconfig.Process("", &cfg.Server); err != nil {
		return nil, fmt.Errorf("failed to load server config: %w", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	// Validate database URL
	if c.Database.URL == "" {
		return fmt.Errorf("DATABASE_URL is required")
	}

	// Validate JWT secrets
	if c.JWT.AccessSecret == "" {
		return fmt.Errorf("JWT_ACCESS_SECRET is required")
	}
	if c.JWT.RefreshSecret == "" {
		return fmt.Errorf("JWT_REFRESH_SECRET is required")
	}

	if len(c.JWT.AccessSecret) < 32 {
		return fmt.Errorf("JWT_ACCESS_SECRET must be at least 32 characters")
	}
	if len(c.JWT.RefreshSecret) < 32 {
		return fmt.Errorf("JWT_REFRESH_SECRET must be at least 32 characters")
	}

	// Ensure secrets are different for security
	if c.JWT.AccessSecret == c.JWT.RefreshSecret {
		return fmt.Errorf("JWT_ACCESS_SECRET and JWT_REFRESH_SECRET must be different")
	}

	// Validate SMTP if provider is smtp
	if c.Email.Provider == "smtp" {
		if c.Email.SMTP.Host == "" {
			return fmt.Errorf("SMTP_HOST is required when EMAIL_PROVIDER is smtp")
		}
		if c.Email.SMTP.Username == "" {
			return fmt.Errorf("SMTP_USER is required when EMAIL_PROVIDER is smtp")
		}
		if c.Email.SMTP.Password == "" {
			return fmt.Errorf("SMTP_PASS is required when EMAIL_PROVIDER is smtp")
		}
	}

	// Validate environment
	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
		"test":        true,
	}
	if !validEnvs[c.App.Environment] {
		return fmt.Errorf("invalid APP_ENV: %s (must be development, staging, production, or test)", c.App.Environment)
	}

	return nil
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// GetServerAddress returns the server address in format :port
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf(":%d", c.Server.Port)
}
