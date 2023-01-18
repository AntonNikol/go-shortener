package server

import (
	"github.com/AntonNikol/go-shortener/internal/app/handlers"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/file"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/inmemory"
	"github.com/AntonNikol/go-shortener/internal/config"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

var repo repositories.Repository

func Run(cfg *config.Config) {
	// Определяем какой репозиторий будет использоваться - память или файл
	if cfg.FileStoragePath != "" {
		repo = repositories.Repository(file.New(cfg.FileStoragePath))
	} else {
		repo = repositories.Repository(inmemory.New())
	}

	h := handlers.New(cfg.BaseURL, repo)

	// Routes
	e := echo.New()
	e.POST("/", h.CreateItem)
	e.GET("/:id", h.GetItem)
	e.POST("api/shorten", h.CreateItemJSON)

	log.Printf("Сервер запущен на адресе: %s", cfg.ServerAddress)

	// Start server
	s := http.Server{
		Addr: cfg.ServerAddress,
	}
	e.Logger.Fatal(e.StartServer(&s))
}
