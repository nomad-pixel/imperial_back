package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateVerificationCode генерирует случайный числовой код заданной длины
func GenerateVerificationCode(length int) (string, error) {
	if length <= 0 {
		length = 6 // По умолчанию 6-значный код
	}

	const digits = "0123456789"
	code := make([]byte, length)

	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		code[i] = digits[num.Int64()]
	}

	return string(code), nil
}
