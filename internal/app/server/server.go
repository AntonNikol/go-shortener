package server

import (
	"context"
	"github.com/AntonNikol/go-shortener/internal/app/handlers"
	"github.com/AntonNikol/go-shortener/internal/app/middlewares"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"strings"
)

func Run(ctx context.Context, cfg *config.Config, repo repositories.Repository) {

	h := handlers.New(cfg.BaseURL, repo, cfg.DBDSN)

	// Routes
	e := echo.New()
	e.Use(carryContextMiddleware(ctx))

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

	e.Use(middlewares.CookieMiddleware)
	e.POST("/", h.CreateItem)
	e.POST("api/shorten", h.CreateItemJSON)
	e.POST("api/shorten/batch", h.CreateItemsList)
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

func carryContextMiddleware(ctx context.Context) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
