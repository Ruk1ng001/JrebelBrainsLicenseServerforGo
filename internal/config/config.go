package config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type Config struct {
	Server  ServerConfig
	License LicenseConfig
	Web     WebConfig
}

type ServerConfig struct {
	Port            string
	Host            string
	ServerVersion   string
	ServerGUID      string
	ProtocolVersion string
}

type LicenseConfig struct {
	PrivateKeyPath     string
	ServerRandomness   string
	OfflineDays        int
	ProlongationPeriod string
}

type WebConfig struct {
	DockerImage string
	BaseURL     string
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:            getEnv("PORT", "8081"),
			Host:            getEnv("HOST", "0.0.0.0"),
			ServerVersion:   getEnv("SERVER_VERSION", "3.2.4"),
			ServerGUID:      getEnv("SERVER_GUID", "a1b4aea8-b031-4302-b602-670a990272cb"),
			ProtocolVersion: getEnv("PROTOCOL_VERSION", "1.1"),
		},
		License: LicenseConfig{
			PrivateKeyPath:     getEnv("PRIVATE_KEY_PATH", ""),
			ServerRandomness:   getEnv("SERVER_RANDOMNESS", "H2ulzLlh7E0="),
			OfflineDays:        180,
			ProlongationPeriod: "607875500",
		},
		Web: WebConfig{
			DockerImage: getEnv("DOCKER_IMAGE", "ruk1ng001/jrebel-license-server"),
			BaseURL:     getEnv("BASE_URL", ""), // 留空则自动检测
		},
	}

	return config, nil
}

// LoadPrivateKey 加载私钥
func LoadPrivateKey(keyPath string) (*rsa.PrivateKey, error) {
	// 如果没有提供路径，使用硬编码的密钥（仅用于演示）
	var keyData []byte
	var err error

	if keyPath != "" {
		keyData, err = os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}
	} else {
		// 使用原始的硬编码密钥（仅用于演示）
		keyData = []byte(defaultPrivateKey)
	}

	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	var privateKey interface{}

	// 尝试解析 PKCS8 格式
	privateKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// 如果失败，尝试解析 PKCS1 格式
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
	}

	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}

	return rsaKey, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 默认私钥
const defaultPrivateKey = `-----BEGIN PRIVATE KEY-----
MIICXAIBAAKBgQDQ93CP6SjEneDizCF1P/MaBGf582voNNFcu8oMhgdTZ/N6qa6O
7XJDr1FSCyaDdKSsPCdxPK7Y4Usq/fOPas2kCgYcRS/iebrtPEFZ/7TLfk39HLuT
Ejzo0/CNvjVsgWeh9BYznFaxFDLx7fLKqCQ6w1OKScnsdqwjpaXwXqiulwIDAQAB
AoGATOQvvBSMVsTNQkbgrNcqKdGjPNrwQtJkk13aO/95ZJxkgCc9vwPqPrOdFbZa
ppZeHa5IyScOI2nLEfe+DnC7V80K2dBtaIQjOeZQt5HoTRG4EHQaWoDh27BWuJoi
p5WMrOd+1qfkOtZoRjNcHl86LIAh/+3vxYyebkug4UHNGPkCQQD+N4ZUkhKNQW7m
pxX6eecitmOdN7Yt0YH9UmxPiW1LyCEbLwduMR2tfyGfrbZALiGzlKJize38shGC
1qYSMvZFAkEA0m6psWWiTUWtaOKMxkTkcUdigalZ9xFSEl6jXFB94AD+dlPS3J5g
NzTEmbPLc14VIWJFkO+UOrpl77w5uF2dKwJAaMpslhnsicvKMkv31FtBut5iK6GW
eEafhdPfD94/bnidpP362yJl8Gmya4cI1GXvwH3pfj8S9hJVA5EFvgTB3QJBAJP1
O1uAGp46X7Nfl5vQ1M7RYnHIoXkWtJ417Kb78YWPLVwFlD2LHhuy/okT4fk8LZ9L
eZ5u1cp1RTdLIUqAiAECQC46OwOm87L35yaVfpUIjqg/1gsNwNsj8HvtXdF/9d30
JIM3GwdytCvNRLqP35Ciogb9AO8ke8L6zY83nxPbClM=
-----END PRIVATE KEY-----`
