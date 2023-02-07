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

var IntServErr = "Internal Server Error"

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
		return echo.NewHTTPError(http.StatusInternalServerError, IntServErr)
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
		return echo.NewHTTPError(http.StatusInternalServerError, IntServErr)
	}

	item := models.Item{
		FullURL:  string(body),
		ShortURL: h.baseURL + "/",
		UserID:   user.Value,
	}

	item, err = h.repository.AddItem(c.Request().Context(), item)
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
		return echo.NewHTTPError(http.StatusInternalServerError, IntServErr)
	}

	item := models.Item{}
	if err := c.Bind(&item); err != nil {
		log.Printf("handler CreateItemJSON json parsing error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "JSON parsing error")
	}

	item.ShortURL = h.baseURL + "/"
	item.UserID = user.Value

	item, err = h.repository.AddItem(c.Request().Context(), item)
	if err != nil {
		if errors.Is(err, postgres.ErrUniqueViolation) {
			log.Printf("ErrUniqueViolation, item %v", item)
			return c.JSON(http.StatusConflict, struct {
				Result string `json:"result"`
			}{Result: h.baseURL + "/" + item.ID})
		}
		log.Printf("CreateItemJSON err: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, IntServErr)
	}

	return c.JSON(http.StatusCreated, struct {
		Result string `json:"result"`
	}{Result: item.ShortURL + item.ID})
}

func (h Handlers) GetItem(c echo.Context) error {
	id := c.Param("id")

	item, err := h.repository.GetItemByID(c.Request().Context(), id)
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
		return echo.NewHTTPError(http.StatusInternalServerError, IntServErr)
	}

	items, err := h.repository.GetItemsByUserID(c.Request().Context(), user.Value)
	if err != nil {
		log.Printf("GetItemsByUserID ошибка: %v", err)
		return c.String(http.StatusNoContent, "")
	}
	log.Printf("GetItemsByUserID найдено items: %d", len(items))

	if h.dbDSN != "" {
		var result []models.ItemResponse
		for _, v := range items {
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
		return echo.NewHTTPError(http.StatusInternalServerError, IntServErr)
	}

	return c.String(http.StatusOK, "")
}

func (h Handlers) CreateItemsList(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, IntServErr)
	}
	if len(body) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Request body is required")
	}

	var itemsRequest []models.ItemList

	err = json.Unmarshal(body, &itemsRequest)
	if err != nil {
		log.Printf("Ошибка парсинга json %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, IntServErr)
	}

	user, err := c.Cookie("user_id")
	if err != nil {
		log.Printf("CreateItemJSON не удалось прочитать куки %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, IntServErr)
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

	result, err := h.repository.AddItemsList(c.Request().Context(), items)
	if err != nil {
		log.Printf("CreateItemsList unable use repository AddItemsList %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, IntServErr)
	}

	var response []models.ItemList
	for k, v := range result {
		r := models.ItemList{ID: k, ShortURL: h.baseURL + "/" + v.ID}
		response = append(response, r)
	}
	return c.JSON(http.StatusCreated, response)

}
