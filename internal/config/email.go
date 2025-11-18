package config

import "os"

type EmailConfig struct {
	Provider string
	SMTP     SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func LoadEmailConfig() EmailConfig {

	return EmailConfig{
		Provider: getEnv("EMAIL_PROVIDER", "console"),
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnv("SMTP_PORT", "587"),
			Username: getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASS", ""),
			From:     getEnv("SMTP_FROM", "noreply@imperial.com"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
