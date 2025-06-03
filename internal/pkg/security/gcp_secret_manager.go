package security

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type GCPSecretManagerClient struct {
	client *secretmanager.Client
}

func NewGCPSecretManagerClient(ctx context.Context) (GCPSecretManagerClient, error) {
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return GCPSecretManagerClient{}, err
	}

	return GCPSecretManagerClient{client: client}, nil
}

func (g GCPSecretManagerClient) GetPrivateKey(ctx context.Context, resourceName string) (*rsa.PrivateKey, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: resourceName,
	}

	result, err := g.client.AccessSecretVersion(ctx, req)
	if err != nil {
		return nil, err
	}

	privateKeyPEM := string(result.Payload.Data)
	privateKey, err := g.parseRSAPrivateKey([]byte(privateKeyPEM))
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (g GCPSecretManagerClient) Close() error {
	return g.client.Close()
}

func (g GCPSecretManagerClient) parseRSAPrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing RSA private key")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
