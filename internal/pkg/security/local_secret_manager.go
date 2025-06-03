package security

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

type LocalSecretManager struct {}

func NewLocalSecretManager() (LocalSecretManager, error) {
	return LocalSecretManager{}, nil
}

func (l LocalSecretManager) GetPrivateKey(_ context.Context, _ string) (*rsa.PrivateKey, error) {
	pemBytes, err := os.ReadFile("./private.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %w", err)
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing RSA private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	return privateKey, nil
}

func (l LocalSecretManager) Close() error {
	return nil
}
