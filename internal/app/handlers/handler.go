package handlers

import (
	"encoding/json"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

type Handlers struct {
	baseURL    string
	repository repositories.Repository
}

func New(baseURL string, repository repositories.Repository) *Handlers {
	return &Handlers{
		baseURL:    baseURL,
		repository: repository,
	}
}

func (h Handlers) CreateItem(c echo.Context) error {
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

	randomString := h.generateUniqueItemID("")
	item := models.Item{
		FullURL:  string(body),
		ShortURL: h.baseURL + "/" + randomString,
		ID:       randomString,
	}

	item, err = h.repository.AddItem(item)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())

	}
	return c.String(http.StatusCreated, item.ShortURL)
}

func (h Handlers) CreateItemJSON(c echo.Context) error {
	randomString := h.generateUniqueItemID("")
	item := models.Item{}

	if err := c.Bind(&item); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга json "+err.Error())
	}
	item.ShortURL = h.baseURL + "/" + randomString
	item.ID = randomString

	item, err := h.repository.AddItem(item)
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

func (h Handlers) GetItem(c echo.Context) error {
	id := c.Param("id")

	item, err := h.repository.GetItemByID(id)
	if err != nil {
		return c.String(http.StatusNotFound, "Ссылка не найдена")
	}

	c.Response().Header().Set("Location", item.FullURL)
	return c.String(http.StatusTemporaryRedirect, "")
}

// получение рандомного id
func (h Handlers) generateUniqueItemID(id string) string {
	randomInt := rand.Intn(999999)
	randomString := strconv.Itoa(randomInt)

	log.Printf("generateUniqueItemID Получение рандомного id: %s", id)
	exists := h.checkItemExist(randomString)
	log.Printf("generateUniqueItemID exists id: %v", exists)

	if randomString != id && !exists {
		return randomString
	}

	return h.generateUniqueItemID(randomString)
}

// проверка есть ли в файле item с таким id
func (h Handlers) checkItemExist(id string) bool {

	log.Printf("checkItemExist проверка на существование item c id: %s", id)

	item, err := h.repository.GetItemByID(id)
	log.Printf("checkItemExist item: %v, err %v", item, err)

	return err == nil
}
