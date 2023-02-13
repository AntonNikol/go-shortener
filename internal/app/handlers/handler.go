package handlers

import (
	"encoding/json"
	"errors"
	"github.com/AntonNikol/go-shortener/internal/app/models"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/labstack/echo/v4"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Handlers struct {
	baseURL    string
	repository repositories.Repository
}

const (
	IntServErr = "Internal Server Error"
)

func New(baseURL string, repository repositories.Repository) *Handlers {
	return &Handlers{
		baseURL:    baseURL,
		repository: repository,
	}
}

func (h Handlers) CreateItemHandler(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, IntServErr)
	}

	if len(body) == 0 {
		return c.String(http.StatusBadRequest, "Request body is required")
	}

	_, err = url.ParseRequestURI(string(body))
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid url")
	}

	log.Printf("CreateItemHandle2")
	user, err := c.Cookie("user_id")
	if err != nil {
		return c.String(http.StatusInternalServerError, IntServErr)
	}
	userID := user.Value

	item := models.Item{
		FullURL: string(body),
		UserID:  userID,
	}

	item, err = h.repository.AddItem(c.Request().Context(), item)
	if err != nil && !errors.Is(err, repositories.ErrAlreadyExists) {
		// а вот пятисотки логгировать как раз надо
		log.Printf("unable to add item %v in repo: %v", item, err)
		return c.String(http.StatusInternalServerError, h.baseURL+"/"+item.ID)
	}
	if errors.Is(err, repositories.ErrAlreadyExists) {
		// нам незачем логгировать ошибки 4хх - иначе тогда любой клиент сможет хранилку логов задудосить
		return c.String(http.StatusConflict, h.baseURL+"/"+item.ID)
	}

	return c.String(http.StatusCreated, h.baseURL+"/"+item.ID)
}

func (h Handlers) CreateItemJSONHandler(c echo.Context) error {
	user, err := c.Cookie("user_id")
	if err != nil {
		log.Printf("CreateItemJSON не удалось прочитать куки %v", err)
		return c.String(http.StatusInternalServerError, IntServErr)
	}
	userID := user.Value

	var item models.Item
	if err := c.Bind(&item); err != nil {
		return c.String(http.StatusBadRequest, "JSON parsing error")
	}
	item.UserID = userID
	item, err = h.repository.AddItem(c.Request().Context(), item)
	if err != nil {
		if errors.Is(err, repositories.ErrAlreadyExists) {
			return c.JSON(http.StatusConflict, struct {
				Result string `json:"result"`
			}{Result: h.baseURL + "/" + item.ID})
		}
		log.Printf("CreateItemJSON err: %v", err)
		return c.String(http.StatusInternalServerError, IntServErr)
	}

	return c.JSON(http.StatusCreated, struct {
		Result string `json:"result"`
	}{Result: h.baseURL + "/" + item.ID})
}

func (h Handlers) GetItemHandler(c echo.Context) error {
	id := c.Param("id")

	item, err := h.repository.GetItemByID(c.Request().Context(), id)
	if err != nil && !errors.Is(err, repositories.ErrNotFound) {
		log.Printf("unable to get item %v from repo: %v", id, err)
		return c.String(http.StatusInternalServerError, h.baseURL+"/"+id)
	}
	if errors.Is(err, repositories.ErrNotFound) {
		return c.String(http.StatusNotFound, "Ссылка не найдена")
	}

	c.Response().Header().Set("Location", item.FullURL)
	return c.String(http.StatusTemporaryRedirect, "")
}

func (h Handlers) GetItemsByUserIDHandler(c echo.Context) error {
	user, err := c.Cookie("user_id")
	if err != nil {
		return c.String(http.StatusInternalServerError, IntServErr)
	}
	userID := user.Value

	items, err := h.repository.GetItemsByUserID(c.Request().Context(), userID)
	if err != nil {
		log.Printf("GetItemsByUserID ошибка: %v", err)
		return c.String(http.StatusNoContent, "")
	}
	log.Printf("GetItemsByUserID найдено items: %d", len(items))

	var result []models.ItemResponse
	for _, v := range items {
		log.Printf("Подстановка v.ShortURL было: %s", v.ShortURL)

		v.ShortURL = h.baseURL + "/" + v.ID
		result = append(result, v)
	}

	return c.JSON(http.StatusOK, result)
}

func (h Handlers) DBPingHandler(c echo.Context) error {
	err := h.repository.Ping(c.Request().Context())
	if err != nil {
		log.Println("err ping")
		return c.String(http.StatusInternalServerError, IntServErr)
	}

	return c.String(http.StatusOK, "")
}

func (h Handlers) CreateItemsListHandler(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, IntServErr)
	}
	if len(body) == 0 {
		return c.String(http.StatusBadRequest, "Request body is required")
	}

	var itemsRequest []models.ItemList

	err = json.Unmarshal(body, &itemsRequest)
	if err != nil {
		log.Printf("Ошибка парсинга json %v", err)
		return c.String(http.StatusBadRequest, IntServErr)
	}

	user, err := c.Cookie("user_id")
	if err != nil {
		return c.String(http.StatusInternalServerError, IntServErr)
	}
	userID := user.Value
	// Собираем мапу айтемсов
	items := make(map[string]models.Item)
	for _, v := range itemsRequest {
		item := models.Item{
			FullURL: v.OriginalURL,
			UserID:  userID,
		}
		items[v.ID] = item
	}

	result, err := h.repository.AddItemsList(c.Request().Context(), items)
	if err != nil {
		log.Printf("CreateItemsList unable use repository AddItemsList %v", err)
		return c.String(http.StatusInternalServerError, IntServErr)
	}
	log.Printf("получен result %+v", result)

	var response []models.ItemList
	for k, v := range result {
		r := models.ItemList{ID: k, ShortURL: h.baseURL + "/" + v.ID}
		response = append(response, r)
	}
	return c.JSON(http.StatusCreated, response)
}
