package utils

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRandomBase64 生成随机Base64字符串
func GenerateRandomBase64(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}
