package signer

import (
	"crypto/hmac"
	"crypto/sha256"
)

func Sign(b []byte, signingKey []byte) []byte {
	h := hmac.New(sha256.New, signingKey)
	return h.Sum(b)
}

func Verify(b []byte, signature []byte, signingKey []byte) bool {
	h := hmac.New(sha256.New, signingKey)
	expectedSignature := h.Sum(b)
	return hmac.Equal(signature, expectedSignature)
}
