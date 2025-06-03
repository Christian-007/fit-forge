package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

func main() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048) // 2048-bit RSA key
	if err != nil {
		log.Fatalf("failed to generate RSA Private Key: %s", err.Error())
	}

	err = savePrivateKey("private.pem", privateKey)
	if err != nil {
		log.Fatalf("failed to save the private key to 'private.pem' file: %s", err.Error())
	}
	log.Print("Successfully saved private key to 'private.pem'")

	jwksJson, err := generateJwks(privateKey)
	if err != nil {
		log.Fatalf("failed to generate JWKS JSON: %s", err.Error())
	}

	err = os.WriteFile("jwks.json", jwksJson, 0644)
	if err != nil {
		log.Fatalf("failed to save to 'jwks.json': %s", err.Error())
	}

	log.Print("Successfully saved JWKS to 'jwks.json'")
}

func savePrivateKey(filename string, key *rsa.PrivateKey) error {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(key)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	return os.WriteFile(filename, pem.EncodeToMemory(pemBlock), 0600)
}

func generateJwks(privateKey *rsa.PrivateKey) ([]byte, error) {
	publicKeyJwk, err := jwk.PublicKeyOf(privateKey.Public())
	if err != nil {
		return nil, fmt.Errorf("failed to create JWK from public key: %s", err.Error())
	}

	publicKeyJwk.Set(jwk.KeyIDKey, os.Getenv("JWK_KEY_ID"))
	publicKeyJwk.Set(jwk.AlgorithmKey, jwt.SigningMethodRS256.Alg()) // match with JWT signing algorithm
	publicKeyJwk.Set(jwk.KeyUsageKey, jwk.ForSignature)

	set := jwk.NewSet()
	set.AddKey(publicKeyJwk)

	// Marshal JWKS with indentation
	buf, err := json.MarshalIndent(set, "", "  ")
	if err != nil {
		return nil, err
	}
	return buf, nil
}
