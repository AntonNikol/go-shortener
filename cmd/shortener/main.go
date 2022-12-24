package main

import (
	"github.com/labstack/echo/v4"
	"io"
	"math/rand"
	"net/http"
	"strconv"
)

var items []Item

var host = "http://localhost:8080"

type Item struct {
	FullURL  string `json:"full_url"`
	ShortURL string `json:"short_url"`
	Id       string
}

func main() {
	e := echo.New() // Routes
	e.GET("/:id", getItem)
	e.POST("/", createItem)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

func createItem(c echo.Context) error {
	defer c.Request().Body.Close()

	if c.Request().Body == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "body обязательно")
	}

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	randomString := strconv.Itoa(rand.Int())
	randomString = randomString[:6]

	item := Item{
		FullURL:  string(body),
		ShortURL: host + "/" + randomString,
		Id:       randomString,
	}
	items = append(items, item)

	return c.String(http.StatusCreated, item.ShortURL)
}

func getItem(c echo.Context) error {
	id := c.Param("id")

	for _, item := range items {
		if item.Id == id {
			c.Response().Header().Set("Location", item.FullURL)

			return c.String(http.StatusTemporaryRedirect, item.FullURL)
		}
	}
	return c.String(404, "Ссылка не найдена")
}

//TODO:
// проверка что body не пустой
// перенести хэндлеры
// сервер
// storage implements interface
