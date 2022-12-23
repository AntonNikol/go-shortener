package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

var items []Item

var host = "http://localhost:8080/"

type Item struct {
	FullURL  string `json:"full_url"`
	ShortURL string `json:"short_url"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", createItem).Methods("POST")
	r.HandleFunc("/", getItems).Methods("GET")
	r.HandleFunc("/{id}", getItem).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}

// Получение полной ссылки по сокращенной ссылке
func getItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	for _, item := range items {
		fmt.Printf("количество элелементов items: %d \n", len(items))

		fmt.Printf("item.ShortURL = %s, id = %s \n", item.ShortURL, params["id"])
		if item.ShortURL == params["id"] {
			//http.Error(w, "Ссылка НАЙДЕНА", http.StatusNotFound)
			//return

			fmt.Printf("условие проверки выполняется, возвращаем ответ\n")
			w.Header().Set("Location", item.FullURL)
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Write([]byte("Готово! \n" + item.FullURL))
			return

		}
	}

	http.Error(w, "Ссылка не найдена", http.StatusNotFound)
}

// Сокращение ссылки
func createItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	randomString := strconv.Itoa(rand.Int())
	randomString = randomString[:6]

	item := Item{
		FullURL:  string(body),
		ShortURL: randomString,
	}
	items = append(items, item)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(host + item.ShortURL))
}

func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

//TODO:
// проверка что body не пустой
// перенести хэндлеры
// сервер
// storage implements interface
