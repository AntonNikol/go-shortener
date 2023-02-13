package middlewares

import (
	"fmt"
	"github.com/AntonNikol/go-shortener/pkg/generator"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"time"
)

func CookieMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("user_id")

		if err != nil || cookie.Value == "" {
			log.Printf("userCookieMiddleware чтение куки, ошибка %v", err)
			log.Printf("userCookieMiddleware куки пустые. пишем новые")
			userID, err := generator.GenerateRandomID(16)

			fmt.Printf("middleware GenerateRandomID")
			if err != nil {
				log.Printf("userCookieMiddleware generateUserID ошибка %v", err)
				return err
			}
			// Устанавливаем куки в заголовки ответа
			cookie := new(http.Cookie)
			cookie.Name = "user_id"
			cookie.Value = userID
			cookie.Expires = time.Now().Add(24 * time.Hour)
			c.SetCookie(cookie)
			log.Println("userCookieMiddleware куки установлены")

			//Установить куки в заголовки запроса
			c.Request().AddCookie(cookie)
		}
		log.Println("userCookieMiddleware конец мидлвара")

		return next(c)
	}
}
