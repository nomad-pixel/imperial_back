# Configuration Guide

This document describes all configuration options for the Imperial application.

## Overview

The application uses environment variables for configuration, loaded and validated using [envconfig](https://github.com/kelseyhightower/envconfig).

## Quick Start

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Edit `.env` with your settings:
```bash
# Required: Set a strong JWT secret (min 32 characters)
JWT_SECRET=$(openssl rand -base64 32)

# Required: Configure database connection
DATABASE_URL=postgres://user:password@localhost:5432/imperial?sslmode=disable
```

3. Run the application:
```bash
make run
```

## Configuration Sections

### Application Settings

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `APP_NAME` | string | No | `Imperial` | Application name for logging |
| `APP_ENV` | string | No | `development` | Environment: `development`, `staging`, `production`, `test` |
| `APP_DEBUG` | bool | No | `false` | Enable debug mode |

**Example:**
```bash
APP_NAME=Imperial
APP_ENV=production
APP_DEBUG=false
```

### Database Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `DATABASE_URL` | string | **Yes** | - | PostgreSQL connection string |
| `DB_MAX_CONNS` | int | No | `25` | Maximum number of connections in pool |
| `DB_MIN_CONNS` | int | No | `5` | Minimum number of connections in pool |
| `DB_MAX_CONN_LIFETIME` | duration | No | `1h` | Maximum lifetime of a connection |
| `DB_MAX_CONN_IDLE_TIME` | duration | No | `30m` | Maximum idle time before closing connection |

**Example:**
```bash
DATABASE_URL=postgres://imperial:secret@localhost:5432/imperial?sslmode=disable
DB_MAX_CONNS=50
DB_MIN_CONNS=10
DB_MAX_CONN_LIFETIME=2h
DB_MAX_CONN_IDLE_TIME=1h
```

**Duration format:** `1s`, `5m`, `1h`, `24h`, etc.

### JWT Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `JWT_SECRET` | string | **Yes** | - | Secret key for JWT signing (min 32 chars) |
| `JWT_ACCESS_DURATION` | duration | No | `15m` | Access token lifetime |
| `JWT_REFRESH_DURATION` | duration | No | `168h` | Refresh token lifetime (7 days) |

**Example:**
```bash
# Generate a secure secret:
# openssl rand -base64 32
JWT_SECRET=your-generated-secret-key-here-min-32-characters
JWT_ACCESS_DURATION=30m
JWT_REFRESH_DURATION=720h  # 30 days
```

**Security Notes:**
- JWT secret MUST be at least 32 characters
- Never commit the secret to version control
- Use different secrets for each environment
- Rotate secrets periodically

### Email Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `EMAIL_PROVIDER` | string | No | `console` | Email provider: `console` or `smtp` |
| `SMTP_HOST` | string | Conditional* | `smtp.gmail.com` | SMTP server hostname |
| `SMTP_PORT` | int | Conditional* | `587` | SMTP server port |
| `SMTP_USER` | string | Conditional* | - | SMTP username |
| `SMTP_PASS` | string | Conditional* | - | SMTP password |
| `SMTP_FROM` | string | Conditional* | `noreply@imperial.com` | Default sender email |

\* Required only when `EMAIL_PROVIDER=smtp`

**Example - Console (Development):**
```bash
EMAIL_PROVIDER=console
```

**Example - SMTP (Production):**
```bash
EMAIL_PROVIDER=smtp
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-app@gmail.com
SMTP_PASS=your-app-specific-password
SMTP_FROM=noreply@yourdomain.com
```

**Gmail Setup:**
1. Enable 2-factor authentication
2. Generate an App Password at: https://myaccount.google.com/apppasswords
3. Use the App Password as `SMTP_PASS`

### Storage Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `STORAGE_TYPE` | string | No | `local` | Storage type (currently only `local`) |
| `STORAGE_LOCAL_PATH` | string | No | `./uploads` | Local filesystem path for uploads |
| `STORAGE_BASE_URL` | string | No | `http://localhost:8080/uploads` | Base URL for accessing uploaded files |

**Example:**
```bash
STORAGE_TYPE=local
STORAGE_LOCAL_PATH=/var/www/imperial/uploads
STORAGE_BASE_URL=https://cdn.yourdomain.com/uploads
```

**Notes:**
- Ensure the path exists and is writable
- For production, consider using a CDN for `STORAGE_BASE_URL`

### Server Configuration

| Variable | Type | Required | Default | Description |
|----------|------|----------|---------|-------------|
| `SERVER_PORT` | int | No | `8080` | HTTP server port |
| `SERVER_READ_TIMEOUT` | duration | No | `10s` | Maximum duration for reading request |
| `SERVER_WRITE_TIMEOUT` | duration | No | `10s` | Maximum duration for writing response |
| `SERVER_SHUTDOWN_TIMEOUT` | duration | No | `5s` | Maximum duration for graceful shutdown |

**Example:**
```bash
SERVER_PORT=3000
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_SHUTDOWN_TIMEOUT=10s
```

## Environment-Specific Configurations

### Development

```bash
APP_ENV=development
APP_DEBUG=true
EMAIL_PROVIDER=console
DATABASE_URL=postgres://imperial:imperial@localhost:5432/imperial_dev?sslmode=disable
JWT_SECRET=dev-secret-key-change-in-production-min-32-chars
STORAGE_BASE_URL=http://localhost:8080/uploads
```

### Staging

```bash
APP_ENV=staging
APP_DEBUG=false
EMAIL_PROVIDER=smtp
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=staging@yourdomain.com
SMTP_PASS=your-app-password
DATABASE_URL=postgres://user:pass@staging-db:5432/imperial_staging?sslmode=require
JWT_SECRET=$(vault read -field=value secret/staging/jwt-secret)
STORAGE_BASE_URL=https://staging-cdn.yourdomain.com/uploads
```

### Production

```bash
APP_ENV=production
APP_DEBUG=false
EMAIL_PROVIDER=smtp
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=your-sendgrid-api-key
DATABASE_URL=postgres://user:pass@prod-db:5432/imperial?sslmode=require
DB_MAX_CONNS=100
JWT_SECRET=$(vault read -field=value secret/production/jwt-secret)
STORAGE_TYPE=local
STORAGE_LOCAL_PATH=/mnt/storage/uploads
STORAGE_BASE_URL=https://cdn.yourdomain.com/uploads
SERVER_PORT=8080
```

## Validation

The application validates configuration on startup. Common validation errors:

### Missing Required Variables

```
Error: failed to load config: required key DATABASE_URL missing value
```

**Solution:** Set the required environment variable.

### Invalid JWT Secret

```
Error: invalid config: JWT_SECRET must be at least 32 characters
```

**Solution:** Generate a longer secret:
```bash
openssl rand -base64 48
```

### Invalid Environment

```
Error: invalid config: invalid APP_ENV: dev (must be development, staging, production, or test)
```

**Solution:** Use one of the valid environment values.

### SMTP Configuration Missing

```
Error: invalid config: SMTP_HOST is required when EMAIL_PROVIDER is smtp
```

**Solution:** Either set `EMAIL_PROVIDER=console` or provide all SMTP settings.

## Loading Configuration in Code

The configuration is automatically loaded and validated during application startup:

```go
// In internal/di/providers.go
func ProvideConfig() (*config.Config, error) {
    cfg, err := config.Load()
    if err != nil {
        return nil, err
    }

    if err := cfg.Validate(); err != nil {
        return nil, err
    }

    return cfg, nil
}
```

### Accessing Config Values

Config is injected into providers via Wire:

```go
func ProvideDatabase(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
    // Use cfg.Database.URL, cfg.Database.MaxConns, etc.
}

func ProvideEmailService(cfg *config.Config) (ports.EmailService, error) {
    if cfg.Email.Provider == "smtp" {
        // Use cfg.Email.SMTP.*
    }
}
```

## Best Practices

1. **Never commit secrets**
   - Add `.env` to `.gitignore`
   - Use `.env.example` for documentation
   - Use secret management tools (Vault, AWS Secrets Manager, etc.) in production

2. **Use strong secrets**
   ```bash
   # Generate 32+ character secrets
   openssl rand -base64 32
   ```

3. **Different secrets per environment**
   - Development: Simple, documented in `.env.example`
   - Staging: Rotated monthly
   - Production: Rotated weekly, stored in secret manager

4. **Validate early**
   - Application validates config on startup
   - Fail fast if configuration is invalid
   - Provides clear error messages

5. **Document all variables**
   - Keep `.env.example` up to date
   - Add comments for complex settings
   - Document valid values and formats

6. **Use duration strings**
   ```bash
   # Good
   JWT_ACCESS_DURATION=15m

   # Bad - no unit
   JWT_ACCESS_DURATION=900
   ```

## Troubleshooting

### Application won't start

1. Check all required variables are set:
   ```bash
   DATABASE_URL
   JWT_SECRET
   ```

2. Validate `.env` syntax:
   ```bash
   # No spaces around =
   JWT_SECRET=value  # Good
   JWT_SECRET = value  # Bad
   ```

3. Check database connectivity:
   ```bash
   psql $DATABASE_URL -c "SELECT 1;"
   ```

### Email not sending

1. If using console mode, check logs
2. If using SMTP:
   - Verify credentials
   - Check firewall/network access to SMTP server
   - Test with telnet: `telnet smtp.gmail.com 587`

### File uploads failing

1. Check storage path exists:
   ```bash
   ls -la $STORAGE_LOCAL_PATH
   ```

2. Verify write permissions:
   ```bash
   touch $STORAGE_LOCAL_PATH/test.txt
   ```

3. Check disk space:
   ```bash
   df -h $STORAGE_LOCAL_PATH
   ```

## Migration from Old Config

If you're migrating from the old configuration system:

**Old (email.go):**
```go
emailConfig := config.LoadEmailConfig()
// Manual env reading with getEnv()
```

**New (config.go):**
```go
cfg, err := config.Load()
// Automatic loading with validation
// Access via cfg.Email.*
```

**Environment variables remain the same**, no changes needed to your `.env` file!

## Reference

- [envconfig documentation](https://github.com/kelseyhightower/envconfig)
- [PostgreSQL connection strings](https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING)
- [Go time.Duration format](https://pkg.go.dev/time#ParseDuration)
