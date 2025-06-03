package security

import (
	"context"
	"crypto/rsa"
)

type SecretManageProvider interface {
	GetPrivateKey(ctx context.Context, resourceName string) (*rsa.PrivateKey, error)
	Close() error
}
