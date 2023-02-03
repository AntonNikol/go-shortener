package config

import (
	"flag"
	"os"
)

// можно использовать библиотеку env config для парсинга флагов

// Пример запуска сервера
// go run cmd/shortener/main.go -a=localhost:8008 -b=http://localhost:8008 -f=items_test.txt -d
//go run cmd/shortener/main.go -a=localhost:8008 -b=http://localhost:8008 -d="postgres://postgres:qwerty@localhost:5439/postgres?sslmode=disable"

// с БД
// docker run --name=postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD='qwerty' -p5438:5432 -d --rm postgres
// go run cmd/shortener/main.go -d="postgres://postgres:qwerty@localhost:5439/postgres?sslmode=disable"

type Config struct {
	BaseURL         string
	ServerAddress   string `env:"server_address"`
	FileStoragePath string
	DBDSN           string
}

func Get() *Config {
	serverAddress := os.Getenv("SERVER_ADDRESS")
	baseURL := os.Getenv("BASE_URL")
	fileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	dbDSN := os.Getenv("DATABASE_DSN")

	address := flag.String("a", serverAddress, "server address")
	url := flag.String("b", baseURL, "base url")
	storage := flag.String("f", fileStoragePath, "file storage path")
	db := flag.String("d", dbDSN, "db address")
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
		DBDSN:           *db,
	}
}
