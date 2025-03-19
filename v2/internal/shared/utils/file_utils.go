package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
)

// EnsureDirectoryExists проверяет и создает директорию, если ее нет
func EnsureDirectoryExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// HashFileName хэширует имя файла
func HashFileName(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}
