package server

import (
	"github.com/AntonNikol/go-shortener/internal/app/handlers"
	"github.com/labstack/echo/v4"
)

func Run() {
	e := echo.New() // Routes
	e.GET("/:id", handlers.GetItem)
	e.POST("/", handlers.CreateItem)
	e.POST("api/shorten", handlers.CreateItemJSON)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
