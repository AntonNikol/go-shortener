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

	//// open sql
	//db, err := sql.Open("postgres",
	//	"postgres://postgres:qwerty@localhost:5438/postgres?sslmode=disable")
	//if err != nil {
	//	panic(err)
	//}
	//defer db.Close()
	//// работаем с базой
	//// ...
	//// можем продиагностировать соединение
	//ctx := context.Background()
	//ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	//defer cancel()
	//if err = db.PingContext(ctx); err != nil {
	//	panic(err)
	//}
	//
	//os.Exit(0)
	//

	//if cfg.DBDSN != "" {
	//	log.Printf("Передан database dsn %s", cfg.DBDSN)
	//	ctx, _ := context.WithCancel(context.Background())
	//	//urlExample := "postgres://postgres:qwerty@localhost:5438/postgres"
	//	//urlExample := "postgres://postgres:qwerty@localhost:5438/postgres"
	//	conn, err := pgx.Connect(ctx, cfg.DBDSN)
	//	if err != nil {
	//		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//		os.Exit(1)
	//	}
	//
	//	err = conn.Ping(ctx)
	//	if err != nil {
	//		log.Println("err ping")
	//	}
	//
	//	log.Println("успешный пинг")
	//	//
	//}

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
1 Кажется, что код handler.go уже пергружен, как лучше его оптимизировать и разнести по файлам
2 Может быть стоит называть функции в хэндлере и репозитории разными именами?
Когда в обработчике и где-то внутри вызываемых классов функции имеют одинаковое название
вводя в поиске название функции находится много лишнего
3 Может быть хендлеры в файле handler.go называть с постфиксом handler? DBPingHandler и т.д..
так как функций в классе уже много

*/
