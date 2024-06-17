package crypt

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"

	"github.com/artem-benda/monitor/internal/logger"
	"go.uber.org/zap"
)

// GenerateKeyPair generates a new key pair
func GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		logger.Log.Error("error generating key pair", zap.Error(err))
	}
	return privkey, &privkey.PublicKey
}

// PrivateKeyToBytes private key to bytes
func PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) []byte {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		logger.Log.Error("error marshalling public key", zap.Error(err))
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes
}

// BytesToPrivateKey bytes to private key
func BytesToPrivateKey(priv []byte) *rsa.PrivateKey {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			logger.Log.Error("error decrypting PEM block", zap.Error(err))
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		logger.Log.Error("error parsing PKCS1 private key", zap.Error(err))
	}
	return key
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) *rsa.PublicKey {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			logger.Log.Error("error decrypting PEM block", zap.Error(err))
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		logger.Log.Error("error parsing PKIX public key", zap.Error(err))
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		logger.Log.Error("parsed public key is not rsa public key")
	}
	return key
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	return rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	return rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
}
