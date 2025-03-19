package utils

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"math/big"
)

// Генерация случайного 6-значного пароля
func GenerateRandomPassword() (string, error) {
	const length = 6
	const digits = "0123456789"

	password := make([]byte, length)
	for i := range password {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		password[i] = digits[random.Int64()]
	}

	return string(password), nil
}

// Хеширование пароля
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Проверка пароля
func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
