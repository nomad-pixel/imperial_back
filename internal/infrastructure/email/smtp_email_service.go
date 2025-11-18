package email

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"mime"
	"net/smtp"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
	apperrors "github.com/nomad-pixel/imperial/pkg/errors"
)

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type SMTPEmailService struct {
	config          SMTPConfig
	templateManager *TemplateManager
}

func NewSMTPEmailService(config SMTPConfig) (ports.EmailService, error) {
	tm, err := NewTemplateManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize template manager: %w", err)
	}

	return &SMTPEmailService{
		config:          config,
		templateManager: tm,
	}, nil
}

func (s *SMTPEmailService) SendVerificationCode(ctx context.Context, email, code string) error {
	subject := "Код верификации email"

	htmlBody, err := s.templateManager.Render("verification_code.html", TemplateData{Code: code})
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeInternal, "Ошибка рендеринга HTML шаблона")
	}

	plainTextBody := s.templateManager.GetPlainText(code)

	return s.sendMultipartEmail(email, subject, plainTextBody, htmlBody)
}

func (s *SMTPEmailService) SendPasswordResetCode(ctx context.Context, email, code string) error {
	subject := "Код сброса пароля"

	htmlBody, err := s.templateManager.Render("password_reset.html", TemplateData{Code: code})
	if err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeInternal, "Ошибка рендеринга HTML шаблона")
	}

	plainTextBody := s.templateManager.GetPlainText(code)

	return s.sendMultipartEmail(email, subject, plainTextBody, htmlBody)
}

func (s *SMTPEmailService) sendMultipartEmail(to, subject, plainText, htmlBody string) error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)

	boundaryBytes := make([]byte, 16)
	if _, err := rand.Read(boundaryBytes); err != nil {
		return apperrors.Wrap(err, apperrors.ErrCodeInternal, "Ошибка генерации boundary")
	}
	boundary := "----=_Part_" + hex.EncodeToString(boundaryBytes)

	var message bytes.Buffer
	message.WriteString(fmt.Sprintf("From: %s\r\n", s.config.From))
	message.WriteString(fmt.Sprintf("To: %s\r\n", to))
	message.WriteString(fmt.Sprintf("Subject: %s\r\n", mime.QEncoding.Encode("UTF-8", subject)))
	message.WriteString("MIME-Version: 1.0\r\n")
	message.WriteString(fmt.Sprintf("Content-Type: multipart/alternative; boundary=%s\r\n", boundary))
	message.WriteString("\r\n")

	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	message.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	message.WriteString("\r\n")
	message.WriteString(plainText)
	message.WriteString("\r\n")

	message.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	message.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	message.WriteString("Content-Transfer-Encoding: quoted-printable\r\n")
	message.WriteString("\r\n")
	message.WriteString(htmlBody)
	message.WriteString("\r\n")

	message.WriteString(fmt.Sprintf("--%s--\r\n", boundary))

	log.Printf("📧 Отправка email на %s через %s", to, addr)
	err := smtp.SendMail(addr, auth, s.config.From, []string{to}, message.Bytes())
	if err != nil {
		log.Printf("❌ Ошибка отправки email: %v", err)
		return apperrors.Wrap(err, apperrors.ErrCodeExternal, fmt.Sprintf("Ошибка отправки email через SMTP: %v", err))
	}

	log.Printf("✅ Email успешно отправлен на %s", to)
	return nil
}
