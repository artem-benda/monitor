package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"

	"github.com/artem-benda/monitor/internal/crypt"
)

func mustUnmarshallRSAPrivateKey(base64PrivateKey string) *rsa.PrivateKey {
	if len(base64PrivateKey) == 0 {
		return nil
	}

	r := bytes.NewReader([]byte(base64PrivateKey))

	var rawKey []byte
	_, err := base64.NewDecoder(base64.StdEncoding, r).Read(rawKey)
	if err != nil {
		panic(err)
	}

	key, err := crypt.BytesToPrivateKey(rawKey)
	if err != nil {
		panic(err)
	}

	return key
}
