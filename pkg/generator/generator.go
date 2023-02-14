package generator

import (
	"encoding/hex"
	"log"
	"math/rand"
)

func GenerateRandomID(len int) (string, error) {
	// определяем слайс нужной длины
	b := make([]byte, len)
	// записываем байты в массив b
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("generateRandomID error: %v\n", err)
		return "", err
	}

	return hex.EncodeToString(b), nil
}
