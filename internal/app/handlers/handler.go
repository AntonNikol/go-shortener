package handlers

import (
	"fmt"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/labstack/echo/v4"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

var items []models.Item

var host = "http://localhost:8080"

func CreateItem(c echo.Context) error {
	defer c.Request().Body.Close()

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	if len(body) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Body обязательно")
	}

	_, err = url.ParseRequestURI(string(body))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Невалидный url")

	}
	fmt.Printf("body: %s\n", string(body))

	randomString := strconv.Itoa(rand.Int())
	randomString = randomString[:6]

	item := models.Item{
		FullURL:  string(body),
		ShortURL: host + "/" + randomString,
		ID:       randomString,
	}
	items = append(items, item)

	return c.String(http.StatusCreated, item.ShortURL)
}

func GetItem(c echo.Context) error {
	id := c.Param("id")

	for _, item := range items {
		if item.ID == id {
			c.Response().Header().Set("Location", item.FullURL)

			return c.String(http.StatusTemporaryRedirect, item.FullURL)
		}
	}
	return c.String(404, "Ссылка не найдена")
}
