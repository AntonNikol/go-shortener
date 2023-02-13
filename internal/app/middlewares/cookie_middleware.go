package middlewares

import (
	"fmt"
	"github.com/AntonNikol/go-shortener/pkg/generator"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func CookieMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("user_id")
		if err != nil || cookie.Value == "" {
			userID, err := generator.GenerateRandomID(16)
			if err != nil {
				return fmt.Errorf("CookieMiddleware generate id error :%w", err)
			}

			// Устанавливаем куки в заголовки ответа
			cookie := new(http.Cookie)
			cookie.Name = "user_id"
			cookie.Value = userID
			cookie.Expires = time.Now().Add(24 * time.Hour)
			c.SetCookie(cookie)
			// Устанавливаем куки в заголовки запроса
			c.Request().AddCookie(cookie)
		}
		return next(c)
	}
}
