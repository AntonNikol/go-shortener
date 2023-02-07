package middlewares

import (
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func CookieMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("user_id")

		if err != nil || cookie.Value == "" {
			log.Printf("userCookieMiddleware чтение куки, ошибка %v", err)
			log.Printf("userCookieMiddleware куки пустые. пишем новые")
			userID, err := generateUserID()
			if err != nil {
				log.Printf("userCookieMiddleware generateUserID ошибка %v", err)
				return err
			}
			// Устанавливаем куки в заголовки
			cookie := new(http.Cookie)
			cookie.Name = "user_id"
			cookie.Value = userID
			cookie.Expires = time.Now().Add(24 * time.Hour)
			c.SetCookie(cookie)
			log.Println("userCookieMiddleware куки установлены")

			//Как установить куки в заголовки запроса
			c.Request().AddCookie(cookie)
		}

		//log.Println("userCookieMiddleware конец мидлвара")
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
