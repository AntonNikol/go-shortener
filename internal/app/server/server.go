package server

import (
	"github.com/AntonNikol/go-shortener/internal/app/handlers"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func Run() {
	e := echo.New() // Routes
	e.GET("/:id", handlers.GetItem)
	e.POST("/", handlers.CreateItem)
	e.POST("api/shorten", handlers.CreateItemJSON)

	//export SERVER_ADDRESS=localhost:8080
	//export BASE_URL=http://localhost:8080

	// Start server
	s := http.Server{
		Addr: os.Getenv("SERVER_ADDRESS"),
	}
	e.Logger.Fatal(e.StartServer(&s))

}
