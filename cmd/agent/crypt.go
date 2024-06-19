package main

import (
	"bytes"
	"crypto/rsa"
	"encoding/base64"

	"github.com/artem-benda/monitor/internal/crypt"
)

func mustParsePublicKey(base64PublicKey string) *rsa.PublicKey {
	if len(base64PublicKey) == 0 {
		return nil
	}

	r := bytes.NewReader([]byte(base64PublicKey))

	var rawKey []byte
	_, err := base64.NewDecoder(base64.StdEncoding, r).Read(rawKey)
	if err != nil {
		panic(err)
	}

	key, err := crypt.BytesToPublicKey(rawKey)
	if err != nil {
		panic(err)
	}

	return key
}
