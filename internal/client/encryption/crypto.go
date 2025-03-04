package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// GenerateEncryptionKey создает ключ для шифрования из пароля пользователя.
func GenerateEncryptionKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, 100000, 32, sha256.New)
}

func EncodeData(plainData []byte) string {
	return base64.StdEncoding.EncodeToString(plainData)
}

// DecodeData декодирует данные из Base64.
func DecodeData(encodedData string) ([]byte, error) {
	decodedData, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return nil, fmt.Errorf("ошибка декодирования данных: %w", err)
	}

	return decodedData, nil
}

// EncryptData - Шифрует данные с использованием AES-256-GCM
func EncryptData(plainData []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	encrypted := aesGCM.Seal(nonce, nonce, plainData, nil)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

// DecryptData - Дешифрует данные с использованием AES-256-GCM
func DecryptData(encryptedData string, key []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("неверный размер шифрованных данных")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	decrypted, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
