package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
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

	randomString, err := h.generateUniqueItemID("")
	if err != nil {
		log.Printf("Ошибка генерации item ID %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

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
	randomString, err := h.generateUniqueItemID("")
	if err != nil {
		log.Printf("Ошибка генерации item ID %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка сервера")
	}

	item := models.Item{}

	if err := c.Bind(&item); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "ошибка парсинга json "+err.Error())
	}
	item.ShortURL = h.baseURL + "/" + randomString
	item.ID = randomString

	item, err = h.repository.AddItem(item)
	if err != nil {
		log.Printf("Ошибка записи в файл %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
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
func (h Handlers) generateUniqueItemID(id string) (string, error) {
	randomInt := rand.Intn(999999)
	randomString := strconv.Itoa(randomInt)

	log.Printf("generateUniqueItemID Получение рандомного id: %s", id)
	exist, err := h.checkItemExist(randomString)
	if err != nil {
		return "", fmt.Errorf("unable to check item exist item by id: %w", err)
	}

	log.Printf("generateUniqueItemID exists id: %v", exist)

	if randomString != id && !exist {
		return randomString, nil
	}

	return h.generateUniqueItemID(randomString)
}

// проверка есть ли в файле item с таким id
func (h Handlers) checkItemExist(id string) (bool, error) {
	_, err := h.repository.GetItemByID(id)

	// проверяем что ошибка не пустая и она не нот фаунд
	if err != nil && !errors.Is(err, repositories.ErrNotFound) {
		return false, fmt.Errorf("unable to get item by id: %w", err)
	}
	return !errors.Is(err, repositories.ErrNotFound), nil
}
