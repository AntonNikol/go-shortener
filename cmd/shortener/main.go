package main

import (
	"github.com/AntonNikol/go-shortener/internal/app/server"
	"github.com/AntonNikol/go-shortener/internal/config"
)

func main() {
	cfg := config.Get()
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