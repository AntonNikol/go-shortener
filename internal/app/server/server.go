package server

import (
	"github.com/AntonNikol/go-shortener/internal/app/handlers"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/file"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/inmemory"
	"github.com/AntonNikol/go-shortener/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"strings"
)

var repo repositories.Repository

func Run(cfg *config.Config) {
	// Определяем какой репозиторий будет использоваться - память или файл
	if cfg.FileStoragePath != "" {
		repo = repositories.Repository(file.New(cfg.FileStoragePath))
	} else {
		repo = repositories.Repository(inmemory.New())
	}

	h := handlers.New(cfg.BaseURL, repo, cfg.DbDSN)

	// Routes
	e := echo.New()

	// Если в запросе клиента есть заголовок Accept-Encoding gzip, то используем сжатие
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return !strings.Contains(c.Request().Header.Get("Accept-Encoding"), "gzip")
		},
	}))

	// Если в запросе клиента есть заголовок Content-Encoding gzip, то используем декомпрессию
	e.Use(middleware.DecompressWithConfig(middleware.DecompressConfig{
		Skipper: func(c echo.Context) bool {
			return !strings.Contains(c.Request().Header.Get("Content-Encoding"), "gzip")
		},
	}))

	e.POST("/", h.CreateItem)
	e.POST("api/shorten", h.CreateItemJSON)
	e.GET("/:id", h.GetItem)
	e.GET("/api/user/urls", h.GetItemsByUserID)
	e.GET("/ping", h.DBPing)

	log.Printf("Сервер запущен на адресе: %s", cfg.ServerAddress)

	// Start server
	s := http.Server{
		Addr: cfg.ServerAddress,
	}
	e.Logger.Fatal(e.StartServer(&s))
}
