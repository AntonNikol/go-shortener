package config

import (
	"flag"
	"os"
)

// можно использовать библиотеку env config для парсинга флагов

// Пример запуска сервера
//go run cmd/shortener/main.go -a=localhost:8008 -b=http://localhost:8008 -f=items_test.txt

type Config struct {
	BaseURL         string
	ServerAddress   string `env:"server_address"`
	FileStoragePath string
}

func Get() *Config {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	baseURL := os.Getenv("BASE_URL")
	fileStoragePath := os.Getenv("FILE_STORAGE_PATH")

	address := flag.String("a", serverAddress, "server address")
	url := flag.String("b", baseURL, "base url")
	storage := flag.String("f", fileStoragePath, "file storage path")
	flag.Parse()

	serverAddress = *address
	baseURL = *url

	if serverAddress == "" {
		serverAddress = "0.0.0.0:8080"
	}
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	return &Config{
		BaseURL:         baseURL,
		ServerAddress:   serverAddress,
		FileStoragePath: *storage,
	}
}
