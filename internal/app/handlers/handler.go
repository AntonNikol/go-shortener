package handlers

import (
	"fmt"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/inmemory"
	"github.com/AntonNikol/go-shortener/internal/app/storage/memory"
	"github.com/labstack/echo/v4"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

var items []models.Item

var host = "http://localhost:8080"

var repo repositories.RepositoryInterface

func init() {
	db := memory.Storage{}
	repo = repositories.NewRepository(inmemory.New(&db))
}

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

	item, err = repo.AddItem(item)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}
	return c.String(http.StatusCreated, item.ShortURL)
}

func GetItem(c echo.Context) error {
	id := c.Param("id")

	item, err := repo.GetItemByID(id)
	if err != nil {
		return c.String(404, "Ссылка не найдена")
	}

	c.Response().Header().Set("Location", item.FullURL)
	return c.String(http.StatusTemporaryRedirect, "")

}
