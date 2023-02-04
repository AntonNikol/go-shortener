package handlers

import (
	"encoding/hex"
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
	"time"
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
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	if len(body) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Request body is required")
	}

	_, err = url.ParseRequestURI(string(body))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid url")
	}

	randomString, err := h.generateUniqueItemID("")
	if err != nil {
		log.Printf("Ошибка генерации item ID %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}
	//Если в куках передан UserID берем его - иначе генерируем новый
	userID, err := getUserIDFromCookies(c)
	if err != nil {
		userID, err = generateUserID()
		if err != nil {
			log.Printf("ошибка генерации UserID %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
		}
		// Устанавливаем куки в заголовки
		setUserIDInCookies(c, userID)
	}

	/*
		Эту логику нужно как-то отделить в хендлере.
		При работе с БД я возвращаю в качестве ID таблицы, чтобы по нему потом получить запись
		При этом поле shortUrl хранит в себе другой идентификатор
	*/

	item := models.Item{
		FullURL:  string(body),
		ShortURL: h.baseURL + "/" + randomString,
		ID:       randomString,
		UserID:   userID,
	}

	item, err = h.repository.AddItem(item)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// начался хардкод
	if h.dbDSN != "" {

		// если работаем с БД, то возвращаем в качесте id h.baseURL + "/" + ID в таблице

		log.Printf("работаем с БД, возвращаем костыльный ответ. h.dbDSN = %s, item = %s", h.dbDSN, item)
		return c.String(http.StatusCreated, h.baseURL+"/"+item.ID)

	}
	return c.String(http.StatusCreated, item.ShortURL)
}

func (h Handlers) CreateItemJSON(c echo.Context) error {
	randomString, err := h.generateUniqueItemID("")
	if err != nil {
		log.Printf("Ошибка генерации item ID %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	//Если в куках передан UserID берем его - иначе генерируем новый
	userID, err := getUserIDFromCookies(c)
	if err != nil {
		userID, err = generateUserID()
		if err != nil {
			log.Printf("ошибка генерации UserID %v", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
		}
		// Устанавливаем куки в заголовки
		setUserIDInCookies(c, userID)
	}

	item := models.Item{}
	if err := c.Bind(&item); err != nil {
		log.Printf("handler CreateItemJSON json parsing error %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "JSON parsing error")
	}
	item.ShortURL = h.baseURL + "/" + randomString
	item.ID = randomString
	item.UserID = userID

	item, err = h.repository.AddItem(item)
	if err != nil {
		log.Printf("Ошибка записи в файл %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	resultShortURL := item.ShortURL
	// начался хардкод
	if h.dbDSN != "" {
		// если работаем с БД, то возвращаем в качесте id h.baseURL + "/" + ID в таблице
		log.Printf("работаем с БД, возвращаем костыльный ответ. h.dbDSN = %s, item = %s", h.dbDSN, item)
		resultShortURL = h.baseURL + "/" + item.ID

	}
	// кончился

	r, err := json.Marshal(struct {
		Result string `json:"result"`
	}{
		Result: resultShortURL,
	})

	if err != nil {
		log.Printf("Ошибка сериализации json %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
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

func (h Handlers) GetItemsByUserID(c echo.Context) error {
	// Если по юзеру ничего не найдено возвращаем 204
	userID, err := getUserIDFromCookies(c)
	if err != nil {
		return c.String(http.StatusNoContent, "")
	}

	items, err := h.repository.GetItemsByUserID(userID)
	if err != nil {
		log.Printf("GetItemsByUserID ошибка: %v", err)
		return c.String(http.StatusNoContent, "")
	}
	log.Printf("GetItemsByUserID найдено items: %d", len(items))

	return c.JSON(http.StatusOK, items)
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

// Генерация уникального UserID
func generateUserID() (string, error) {
	// определяем слайс нужной длины
	b := make([]byte, 16)
	_, err := rand.Read(b) // записываем байты в массив b
	if err != nil {
		log.Printf("generateUserID error: %v\n", err)
		return "", err
	}

	return hex.EncodeToString(b), nil
}

// Получение UserID из cookie
func getUserIDFromCookies(c echo.Context) (string, error) {
	cookie, err := c.Cookie("user_id")
	if err != nil {
		return "", err
	}
	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)
	return cookie.Value, nil
}

func setUserIDInCookies(c echo.Context, userID string) {
	// Устанавливаем куки в заголовки
	cookie := new(http.Cookie)
	cookie.Name = "user_id"
	cookie.Value = userID
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
}

func (h Handlers) DBPing(c echo.Context) error {

	err := h.repository.Ping(c.Request().Context())
	if err != nil {
		log.Println("err ping")
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
	}

	return c.String(http.StatusOK, "")
}
