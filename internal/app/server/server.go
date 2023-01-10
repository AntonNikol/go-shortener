package server

import (
	"github.com/AntonNikol/go-shortener/internal/app/handlers"
	"github.com/AntonNikol/go-shortener/internal/config"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

func Run(cfg *config.Config) {
	e := echo.New() // Routes
	e.POST("/", handlers.CreateItem)
	e.GET("/:id", handlers.GetItem)
	e.POST("api/shorten", handlers.CreateItemJSON)

	////export SERVER_ADDRESS=localhost:8080
	////export BASE_URL=http://localhost:8080

	log.Printf("Сервер запущен на адресе %s", cfg.ServerAddress)

	// Start server
	s := http.Server{
		Addr: cfg.ServerAddress,
	}
	e.Logger.Fatal(e.StartServer(&s))

}
