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
		repo = repositories.Repository(postgres.New(ctx, cfg.DBDSN))
	}

	server.Run(ctx, cfg, repo)
}

/* TODO вопросы ментору
1 Кажется, что код handler.go уже пергружен, как лучше его оптимизировать и разнести по файлам?
2 Может быть хендлеры в файле handler.go называть с постфиксом handler? DBPingHandler и т.д..
4 Может быть можно как-то все зависимости закинуть в contextWithValue и передавать из main в server.go только его?

//Часть по SQL
1 Как в хендлере определить какая реализация репозитория сейчас используется
и не вызывать код GetItemsByUserID который выполняется внутри условия if h.dbDSN != "" {} так как я не храню short_url в БД
мне приходится его определять перед респонзом
*/

// Планы
// Посмотреть покрытие кода тестами. Покрыть
// Сделать групповое добавление в память и в файл
// Запуск приложения в докер
