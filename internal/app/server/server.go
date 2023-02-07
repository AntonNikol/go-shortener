package server

import (
	"context"
	"encoding/hex"
	"github.com/AntonNikol/go-shortener/internal/app/handlers"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func Run(ctx context.Context, cfg *config.Config, repo repositories.Repository) {

	h := handlers.New(cfg.BaseURL, repo, cfg.DBDSN)

	// Routes
	e := echo.New()
	e.Use(carryContextMiddleware(ctx), userCookieMiddleware)

	// Если в запросе клиента есть заголовок Accept-Encoding gzip, то используем сжатие
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			return !strings.Contains(c.Request().Header.Get("Accept-Encoding"), "gzip")
		},
	}))

	// Если в запросе клиента есть заголовок Content-Encoding gzip, то используем декомпрессию
	e.Use(middleware.DecompressWithConfig(middleware.DecompressConfig{
		Skipper: func(c echo.Context) bool {
			return !strings.Contains(c.Request().Header.Get("Content-Encoding"), "gzip")
		},
	}))

	e.POST("/", h.CreateItem)
	e.POST("api/shorten", h.CreateItemJSON)
	e.POST("api/shorten/batch", h.CreateItemsList)
	e.GET("/:id", h.GetItem)
	e.GET("/api/user/urls", h.GetItemsByUserID)
	e.GET("/ping", h.DBPing)

	log.Printf("Сервер запущен на адресе: %s", cfg.ServerAddress)

	// Start server
	s := http.Server{
		Addr: cfg.ServerAddress,
	}
	e.Logger.Fatal(e.StartServer(&s))
}

func carryContextMiddleware(ctx context.Context) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

func userCookieMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("user_id")
		if err != nil {
			return err
		}
		if cookie.Value == "" {
			userID, err := generateUserID()
			if err != nil {
				return err
			}
			// Устанавливаем куки в заголовки
			cookie := new(http.Cookie)
			cookie.Name = "user_id"
			cookie.Value = userID
			cookie.Expires = time.Now().Add(24 * time.Hour)
			c.SetCookie(cookie)
		}

		return next(c)
	}
}

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

//
//////Если в куках передан UserID берем его - иначе генерируем новый
//userID, err := getUserIDFromCookies(c)
//if err != nil {
//userID, err = generateUserID()
//if err != nil {
//log.Printf("ошибка генерации UserID %v", err)
//return echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
//}
//// Устанавливаем куки в заголовки
//setUserIDInCookies(c, userID)
////}
