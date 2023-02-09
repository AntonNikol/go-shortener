package config

import (
	"flag"
	"os"
)

// можно использовать библиотеку env config для парсинга флагов

// Пример запуска сервера для работы с памятью
//go run cmd/shortener/main.go -a=localhost:8008 -b=http://localhost:8008
// Пример запуска сервера для работы с файлом
//go run cmd/shortener/main.go -a=localhost:8008 -b=http://localhost:8008 -f=items_test.txt -d
// Пример запуска сервера для работы с БД
//go run cmd/shortener/main.go -a=localhost:8008 -b=http://localhost:8008 -d="postgres://postgres:postgres@localhost:54322/postgres?sslmode=disable"

// создание миграции  migrate create -ext sql -dir ./internal/migrations -seq alter_table_short_links
// Запуск БД
//  docker run --name=postgres_practicum -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD='postgres' -p54322:5432 -d postgres

// Откат миграций
//  migrate -source file://./internal/migrations -database "postgresql://postgres:postgres@localhost:54322/postgres?sslmode=disable" down 1

// запуск выбранного теста в нужной итерации
// shortenertest -test.v -test.run=^TestIteration8/TestGzipCompress/shorten$ -binary-path=cmd/shortener/shortener

// запуск автотестов докер
// docker compose -f docker-compose.yml --profile shortener up --remove-orphans --force-recreate -d --build

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
