package email

import (
	"context"
	"fmt"

	"github.com/nomad-pixel/imperial/internal/domain/ports"
)

type ConsoleEmailService struct{}

func NewConsoleEmailService() ports.EmailService {
	return &ConsoleEmailService{}
}

func (s *ConsoleEmailService) SendVerificationCode(ctx context.Context, email, code string) error {
	fmt.Println("=====================================")
	fmt.Println("📧 EMAIL VERIFICATION CODE")
	fmt.Println("=====================================")
	fmt.Printf("To: %s\n", email)
	fmt.Printf("Verification Code: %s\n", code)
	fmt.Println("-------------------------------------")
	fmt.Println("Пожалуйста, используйте этот код для верификации вашего email.")
	fmt.Println("Код действителен в течение 15 минут.")
	fmt.Println("=====================================")
	return nil
}

func (s *ConsoleEmailService) SendPasswordResetCode(ctx context.Context, email, code string) error {
	fmt.Println("=====================================")
	fmt.Println("🔐 PASSWORD RESET CODE")
	fmt.Println("=====================================")
	fmt.Printf("To: %s\n", email)
	fmt.Printf("Reset Code: %s\n", code)
	fmt.Println("-------------------------------------")
	fmt.Println("Пожалуйста, используйте этот код для сброса пароля.")
	fmt.Println("Код действителен в течение 15 минут.")
	fmt.Println("=====================================")
	return nil
}
