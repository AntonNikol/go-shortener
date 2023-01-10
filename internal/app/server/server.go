package server

import (
	"github.com/AntonNikol/go-shortener/internal/app/handlers"
	"github.com/labstack/echo/v4"
	"log"
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

	serverAddress := os.Getenv("SERVER_ADDRESS")
	if serverAddress == "" {
		serverAddress = "0.0.0.0:8080"
	}

	log.Printf("Сервер запущен на адресе %s", serverAddress)
	// Start server
	s := http.Server{
		Addr: serverAddress,
	}
	e.Logger.Fatal(e.StartServer(&s))

}
