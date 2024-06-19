package main

import (
	"encoding/base64"
	"fmt"

	"github.com/artem-benda/monitor/internal/crypt"
)

// Утилита для генерации тестовой пары ключей
func main() {
	// Генерируем пару ключей
	priv, pub, err := crypt.GenerateKeyPair(8192)
	if err != nil {
		panic(err)
	}

	// Подготавливаем данные
	privBytes := crypt.PrivateKeyToBytes(priv)
	pubBytes, err := crypt.PublicKeyToBytes(pub)
	if err != nil {
		panic(err)
	}

	privBase64 := base64.StdEncoding.EncodeToString(privBytes)
	pubBase64 := base64.StdEncoding.EncodeToString(pubBytes)

	fmt.Println("Private key:")
	fmt.Println(privBase64)
	fmt.Println("\nPublic key:")
	fmt.Println(pubBase64)
}
