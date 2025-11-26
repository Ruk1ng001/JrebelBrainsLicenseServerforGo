package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
)

type Signer struct {
	privateKey *rsa.PrivateKey
}

// NewSigner 创建签名器
func NewSigner(privateKey *rsa.PrivateKey) *Signer {
	return &Signer{
		privateKey: privateKey,
	}
}

// SignLeaseData 签名租约数据
func (s *Signer) SignLeaseData(clientRandomness, serverRandomness, guid string, offline bool, validFrom, validUntil string) (string, error) {
	var dataToSign string

	if offline {
		dataToSign = strings.Join([]string{
			clientRandomness,
			serverRandomness,
			guid,
			fmt.Sprintf("%t", offline),
			validFrom,
			validUntil,
		}, ";")
	} else {
		dataToSign = strings.Join([]string{
			clientRandomness,
			serverRandomness,
			guid,
			fmt.Sprintf("%t", offline),
		}, ";")
	}

	return s.Sign([]byte(dataToSign))
}

// Sign 使用SHA1withRSA签名数据
func (s *Signer) Sign(data []byte) (string, error) {
	hashed := sha1.Sum(data)

	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA1, hashed[:])
	if err != nil {
		return "", fmt.Errorf("failed to sign data: %w", err)
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// SignXML 签名XML数据（用于JetBrains产品）
func (s *Signer) SignXML(xmlContent string) (string, error) {
	return s.Sign([]byte(xmlContent))
}
