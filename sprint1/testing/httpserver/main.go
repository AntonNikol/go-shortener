package main

import (
	"log"
	"net/http"
	// важно: путь до пакета handlers в вашем случае может быть другим
	"github.com/AntonNikol/go-shortener/sprint1/testing/httpserver/handlers"
)

func main() {
	http.HandleFunc("/status", handlers.StatusHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
