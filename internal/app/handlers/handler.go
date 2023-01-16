package handlers

import (
	"encoding/json"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/file"
	"github.com/labstack/echo/v4"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

var host = os.Getenv("BASE_URL")
var repo repositories.Repository

func init() {
	repo = repositories.Repository(file.New("items.txt"))
}

func CreateItem(c echo.Context) error {

	// TODO: переделать получение host через конфиг
	host = os.Getenv("BASE_URL")
	if host == "" {
		host = "http://localhost:8080"
	}

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

	//log.Printf("handler CreteItem body: %s\n", string(body))

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

func CreateItemJSON(c echo.Context) error {
	// TODO: переделать получение host через конфиг
	host = os.Getenv("BASE_URL")
	if host == "" {
		host = "http://localhost:8080"
	}

	randomString := getRandomString("")
	item := models.Item{
		ShortURL: host + "/" + randomString,
		ID:       randomString,
	}

	if err := c.Bind(&item); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга json"+err.Error())
	}

	item, err := repo.AddItem(item)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	r, err := json.Marshal(struct {
		Result string `json:"result"`
	}{
		Result: item.ShortURL,
	})

	if err != nil {
		panic(err)
	}

	c.Response().Header().Set("Content-Type", "application/json; charset=UTF-8")
	return c.String(http.StatusCreated, string(r))
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
	randomInt := rand.Intn(999999)
	randomString := strconv.Itoa(randomInt)

	if randomString != id {
		return randomString
	}

	return getRandomString(randomString)
}
