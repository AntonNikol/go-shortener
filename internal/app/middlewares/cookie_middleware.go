package middlewares

import (
	"github.com/AntonNikol/go-shortener/pkg/ctxdata"
	"github.com/AntonNikol/go-shortener/pkg/generator"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const cookieName = "user_id"

func CookieMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, _ := generator.GenerateRandomID(16)
		cookie, err := c.Cookie("user_id")
		if err == nil {
			userID = cookie.Value
		}
		// устанавливаем user_id в контекст
		ctx := ctxdata.SetUserID(c.Request().Context(), userID)
		c.SetRequest(c.Request().WithContext(ctx))
		setCookieResponse(c, userID)

		return next(c)
	}
}

// setCookieResponse Устанавливаем куки в заголовки ответа
func setCookieResponse(c echo.Context, userID string) {
	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = userID
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.SetCookie(cookie)
}
