package main

import (
	"crypto/rsa"
	"encoding/base64"

	"github.com/artem-benda/monitor/internal/crypt"
)

func mustParseRSAPrivateKey(base64PrivateKey string) *rsa.PrivateKey {
	if len(base64PrivateKey) == 0 {
		return nil
	}

	rawKey, err := base64.StdEncoding.DecodeString(base64PrivateKey)
	if err != nil {
		panic(err)
	}

	key, err := crypt.BytesToPrivateKey(rawKey)
	if err != nil {
		panic(err)
	}

	return key
}
