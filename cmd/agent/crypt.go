package main

import (
	"crypto/rsa"
	"encoding/base64"

	"github.com/artem-benda/monitor/internal/crypt"
)

func mustParsePublicKey(base64PublicKey string) *rsa.PublicKey {
	if len(base64PublicKey) == 0 {
		return nil
	}

	rawKey, err := base64.StdEncoding.DecodeString(base64PublicKey)
	if err != nil {
		panic(err)
	}

	key, err := crypt.BytesToPublicKey(rawKey)
	if err != nil {
		panic(err)
	}

	return key
}
