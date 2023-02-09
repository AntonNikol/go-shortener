package generator

import (
	"encoding/hex"
	"log"
	"math/rand"
)

func GenerateRandomID(len int) (string, error) {
	// определяем слайс нужной длины
	b := make([]byte, len)
	_, err := rand.Read(b) // записываем байты в массив b
	if err != nil {
		log.Fatalf("generateRandomID error: %v\n", err)
		return "", err
	}

	return hex.EncodeToString(b), nil
}
