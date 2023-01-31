package main

import (
	"github.com/AntonNikol/go-shortener/internal/app/server"
	"github.com/AntonNikol/go-shortener/internal/config"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
)

//docker run --name=postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD='qwerty' -p5438:5432 -d --rm postgres

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

	//if cfg.DbDSN != "" {
	//	log.Printf("Передан database dsn %s", cfg.DbDSN)
	//	ctx, _ := context.WithCancel(context.Background())
	//	//urlExample := "postgres://postgres:qwerty@localhost:5438/postgres"
	//	//urlExample := "postgres://postgres:qwerty@localhost:5438/postgres"
	//	conn, err := pgx.Connect(ctx, cfg.DbDSN)
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

	server.Run(cfg)
}

/* TODO вопросы ментору
1 Кажется, что код handler.go уже пергружен, как лучше его оптимизировать и разнести по файлам
2 Может быть стоит называть функции в хэндлере и репозитории разными именами?
Когда в обработчике и где-то внутри вызываемых классов функции имеют одинаковое название
вводя в поиске название функции находится много лишнего
3 Может быть хендлеры в файле handler.go называть с постфиксом handler? DBPingHandler и т.д..
так как функций в классе уже много

*/
