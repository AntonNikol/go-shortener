package main

import (
	"github.com/AntonNikol/go-shortener/internal/app/server"
	"github.com/AntonNikol/go-shortener/internal/config"
)

func main() {
	cfg := config.Get()
	server.Run(cfg)
}

//TODO: вопросы ментору
/*
1 Может быть стоит вынести определения репозитория в config.Get()
и в структуре cfg сразу возвращать repository?

2 Правильно ли будет использовать в cfg теги env вместо текущей реализации

3 Не получилось сделать кастомный мидлвар, поэтому сделал через echo

4
*/
