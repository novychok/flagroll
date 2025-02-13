package jwts

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

const (
	publicKeyPath  = "data/keys/public.pem"
	privateKeyPath = "data/keys/private.pem"
	privateKeyBits = 2048
)

type SecretManager struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

func New() (*SecretManager, error) {
	secretManager := &SecretManager{}

	_, err := os.Stat(privateKeyPath)
	if os.IsNotExist(err) {

		protectedKey, err := rsa.GenerateKey(rand.Reader, privateKeyBits)
		if err != nil {
			return nil, err
		}

		publicKey := &protectedKey.PublicKey

		privateBytes := x509.MarshalPKCS1PrivateKey(protectedKey)
		privatePem := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateBytes,
		})

		err = os.WriteFile(privateKeyPath, privatePem, 0600)
		if err != nil {
			return nil, err
		}

		publicBytes := x509.MarshalPKCS1PublicKey(publicKey)
		publicPem := pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicBytes,
		})

		err = os.WriteFile(publicKeyPath, publicPem, 0644)
		if err != nil {
			return nil, err
		}

		secretManager.publicKey = publicKey
		secretManager.privateKey = protectedKey
	} else {
		privatePem, err := os.ReadFile(privateKeyPath)
		if err != nil {
			return nil, err
		}

		privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePem)
		if err != nil {
			return nil, err
		}

		publicPem, err := os.ReadFile(publicKeyPath)
		if err != nil {
			return nil, err
		}

		publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPem)
		if err != nil {
			return nil, err
		}

		secretManager.publicKey = publicKey
		secretManager.privateKey = privateKey
	}

	return secretManager, nil
}

func (s *SecretManager) PublicKey() *rsa.PublicKey {
	return s.publicKey
}

func (s *SecretManager) PrivateKey() *rsa.PrivateKey {
	return s.privateKey
}
