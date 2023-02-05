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
1 Кажется, что код handler.go уже пергружен, как лучше его оптимизировать и разнести по файлам
2 Может быть стоит называть функции в хэндлере и репозитории разными именами?
Когда в обработчике и где-то внутри вызываемых классов функции имеют одинаковое название
вводя в поиске название функции находится много лишнего
3 Может быть хендлеры в файле handler.go называть с постфиксом handler? DBPingHandler и т.д..
так как функций в классе уже много
(!)4 Может быть стоит все зависимости закинуть в contextWithValue и передавать в server.go только его?

//Часть по SQL
5 В sql.Open не передается контекст, это норм?
6 Как правильно тянуть контекст до GetItemById postgres
(!)7 Как в хендлере определить какая реализация репозитория и не вызывать метод generateRandomString для postgres
8 "LastInsertId is not supported by this driver" не отрабатывает в postgres
lastId, err := r.LastInsertId()
		if err != nil {
			return nil, err
		}
*/

// Дела на завтра.
// Инсерт с получением селекта назад
// get item и getItemsByUser не возвращают данные
// избавиться от проверки checkItemExist для запросов в БД
