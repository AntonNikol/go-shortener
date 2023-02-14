package main

import (
	"context"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/file"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/inmemory"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/postgres"
	"github.com/AntonNikol/go-shortener/internal/app/server"
	"github.com/AntonNikol/go-shortener/internal/config"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"log"
)

var repo repositories.Repository

func main() {
	cfg := config.Get()
	ctx := context.Background()
	// Определяем какой репозиторий будет использоваться - память или файл
	if cfg.FileStoragePath != "" {
		repo = repositories.Repository(file.New(cfg.FileStoragePath))
	} else {
		repo = repositories.Repository(inmemory.New())
	}

	if cfg.DBDSN != "" {
		pgs, err := postgres.New(ctx, cfg.DBDSN)
		if err != nil {
			log.Fatal(err)
		}
		repo = repositories.Repository(pgs)
	}
	log.Printf("main go переходим к запуску сервера")

	server.Run(ctx, cfg, repo)
}

/* TODO вопросы ментору
1 Кажется, что код handler.go уже пергружен, как лучше его оптимизировать и разнести по файлам?
*/

// Планы
// Посмотреть покрытие кода тестами. Покрыть
// Сделать групповое добавление в память и в файл
// Запуск приложения в докер
