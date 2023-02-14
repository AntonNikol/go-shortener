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

//var repo repositories.Repository

func Run(ctx context.Context, cfg *config.Config, repo repositories.Repository) {

	h := handlers.New(cfg.BaseURL, repo)

	// Routes
	e := echo.New()
	e.Use(carryContextMiddleware(ctx))
	e.Use(middlewares.CookieMiddleware)

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

	e.POST("/", h.CreateItemHandler)
	e.POST("api/shorten", h.CreateItemJSONHandler)
	e.POST("api/shorten/batch", h.CreateItemsListHandler)
	e.GET("/:id", h.GetItemHandler)
	e.GET("/api/user/urls", h.GetItemsByUserIDHandler)
	e.GET("/ping", h.DBPingHandler)

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
