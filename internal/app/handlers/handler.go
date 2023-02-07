package handlers

import (
	"encoding/json"
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/postgres"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Handlers struct {
	baseURL    string
	repository repositories.Repository
	dbDSN      string
}

func New(baseURL string, repository repositories.Repository, dbDSN string) *Handlers {
	return &Handlers{
		baseURL:    baseURL,
		repository: repository,
		dbDSN:      dbDSN,
	}
}

func (h Handlers) CreateItem(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Printf("CreateItem не удалось прочитать body %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	if len(body) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Request body is required")
	}

	_, err = url.ParseRequestURI(string(body))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid url")
	}

	user, err := c.Cookie("user_id")
	if err != nil {
		log.Printf("CreateItem не удалось прочитать куки %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	item := models.Item{
		FullURL:  string(body),
		ShortURL: h.baseURL + "/",
		UserID:   user.Value,
	}

	item, err = h.repository.AddItem(item)
	if err != nil {
		if errors.Is(err, postgres.ErrUniqueViolation) {
			log.Printf("ErrUniqueViolation, item %v", item)
			return c.String(http.StatusConflict, h.baseURL+"/"+item.ID)
		}
		log.Printf("Handler AddItem, метод repository.AddItem, ошибка %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.String(http.StatusCreated, item.ShortURL)
}

func (h Handlers) CreateItemJSON(c echo.Context) error {
	user, err := c.Cookie("user_id")
	if err != nil {
		log.Printf("CreateItemJSON не удалось прочитать куки %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	item := models.Item{}
	if err := c.Bind(&item); err != nil {
		log.Printf("handler CreateItemJSON json parsing error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "JSON parsing error")
	}

	item.ShortURL = h.baseURL + "/"
	item.UserID = user.Value

	item, err = h.repository.AddItem(item)
	if err != nil {
		if errors.Is(err, postgres.ErrUniqueViolation) {
			log.Printf("ErrUniqueViolation, item %v", item)
			return c.JSON(http.StatusConflict, struct {
				Result string `json:"result"`
			}{Result: h.baseURL + "/" + item.ID})
		}
		log.Printf("CreateItemJSON err: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.JSON(http.StatusCreated, struct {
		Result string `json:"result"`
	}{Result: item.ShortURL + item.ID})
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

func (h Handlers) GetItemsByUserID(c echo.Context) error {
	user, err := c.Cookie("user_id")
	if err != nil {
		log.Printf("GetItemsByUserID не удалось прочитать куки %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	items, err := h.repository.GetItemsByUserID(user.Value)
	if err != nil {
		log.Printf("GetItemsByUserID ошибка: %v", err)
		return c.String(http.StatusNoContent, "")
	}
	log.Printf("GetItemsByUserID найдено items: %d", len(items))

	//TODO: 1
	if h.dbDSN != "" {
		var result []models.ItemResponse
		for _, v := range items {
			log.Printf("Подстановка v.ShortURL было: %s", v.ShortURL)

			v.ShortURL = h.baseURL + "/" + v.ID
			result = append(result, v)
		}

		return c.JSON(http.StatusOK, result)
	}

	return c.JSON(http.StatusOK, items)
}

func (h Handlers) DBPing(c echo.Context) error {

	err := h.repository.Ping(c.Request().Context())
	if err != nil {
		log.Println("err ping")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.String(http.StatusOK, "")
}

func (h Handlers) CreateItemsList(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	if len(body) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Request body is required")
	}

	var itemsRequest []models.ItemList

	err = json.Unmarshal(body, &itemsRequest)
	if err != nil {
		log.Printf("Ошибка парсинга json %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	user, err := c.Cookie("user_id")
	if err != nil {
		log.Printf("CreateItemJSON не удалось прочитать куки %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	// Собираем мапу айтемсов
	items := make(map[string]models.Item)
	for _, v := range itemsRequest {
		item := models.Item{
			FullURL: v.OriginalURL,
			UserID:  user.Value,
		}
		items[v.ID] = item
	}

	result, err := h.repository.AddItemsList(items)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var response []models.ItemList
	for k, v := range result {
		r := models.ItemList{ID: k, ShortURL: h.baseURL + "/" + v.ID}
		response = append(response, r)
	}
	return c.JSON(http.StatusCreated, response)

}
