package handlers

import (
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/inmemory"
	"github.com/AntonNikol/go-shortener/internal/app/storage/memory"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

var host = "http://localhost:8080"

var repo repositories.RepositoryInterface

func init() {
	db := memory.Storage{}
	repo = repositories.RepositoryInterface(inmemory.New(&db))
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

	log.Printf("handler CreteItem body: %s\n", string(body))

	randomString := getRandomString("")
	item := models.Item{
		FullURL:  string(body),
		ShortURL: host + "/" + randomString,
		ID:       randomString,
	}

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
		return c.String(http.StatusNotFound, "Ссылка не найдена")
	}

	c.Response().Header().Set("Location", item.FullURL)
	return c.String(http.StatusTemporaryRedirect, "")
}

func getRandomString(id string) string {
	randomInt := rand.Intn(9000000 - 1000000)
	randomString := strconv.Itoa(randomInt)

	if randomString != id {
		return randomString
	}

	return getRandomString(randomString)
}
