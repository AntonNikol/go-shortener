package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/AntonNikol/go-shortener/internal/app/handlers"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestItem struct {
	FullURL  string
	ShortURL string
}

// Слайс для тестов получения редиректа по сокращенной ссылке
var testItems []TestItem

// Тест сокращения ссылки
func Test_createItem(t *testing.T) {

	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		body string
		want want
	}{
		{
			name: "Обычная ссылка",
			body: "https://practicum.yandex.ru/",
			want: want{
				statusCode:  201,
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			name: "Длинная ссылка",
			body: "https://www.google.com/search?q=goland+%D1%83%D1%80%D0%BE%D0%BA%D0%B8&oq=goland+%D1%83%D1%80%D0%BE%D0%BA%D0%B8&aqs=chrome..69i57j0i10i512.3638j0j15&sourceid=chrome&ie=UTF-8",
			want: want{
				statusCode:  201,
				contentType: "text/plain; charset=UTF-8",
			},
		},
	}

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.body)))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Проверки
			if assert.NoError(t, handlers.CreateItem(c)) {
				require.Equal(t, tt.want.statusCode, rec.Code)
				require.Equal(t, tt.want.contentType, rec.Header().Get("Content-type"))

				// Получаем body ответа
				responseBody := rec.Body.String()
				t.Logf("Ответ сервера %s", responseBody)

				// Проверка, что в ответе url
				_, err := url.ParseRequestURI(responseBody)
				require.NoError(t, err)

				// Проверка слайса items

				//TODO: почему-то len (items) возвращает 0 и далее тесты не падают
				//assert.Equal(t, index+1, len(items))
				fmt.Println(index)

				// Получаем сокращенный url заполняем слайс testItems
				testItems = append(testItems, TestItem{
					FullURL:  tt.body,
					ShortURL: responseBody,
				})
			}
		})
	}
}

// Тест получения полной ссылки по сокращенной
func Test_getItem(t *testing.T) {
	type want struct {
		statusCode int
		response   string
		location   string
		FullURL    string
	}

	tests := []struct {
		name string
		want want
	}{
		{
			name: "Тест получения обычной ссылки",
			want: want{
				statusCode: 307,
				location:   "https://practicum.yandex.ru/",
			},
		},
		{
			name: "Тест получения длинной ссылки",
			want: want{
				statusCode: 307,
				location:   "https://www.google.com/search?q=goland+%D1%83%D1%80%D0%BE%D0%BA%D0%B8&oq=goland+%D1%83%D1%80%D0%BE%D0%BA%D0%B8&aqs=chrome..69i57j0i10i512.3638j0j15&sourceid=chrome&ie=UTF-8",
			},
		},
		{
			name: "Тест запроса по несуществующей ссылке",
			want: want{
				statusCode: 404,
				location:   "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.want.location)))
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			handlers.CreateItem(ctx)
			responseBody := rec.Body.String()

			// Получаем id из ссылки
			split := strings.Split(responseBody, "/")
			splitLen := len(split)
			id := split[splitLen-1]

			req = httptest.NewRequest(http.MethodGet, "/", nil)
			rec = httptest.NewRecorder()

			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(id)

			// Assertions
			if assert.NoError(t, handlers.GetItem(c)) {
				assert.Equal(t, tt.want.statusCode, rec.Code)

				// Если проверяем только то, что при осутствующем id хендлер вернет 404, то завершаем тест
				if tt.name == "Тест запроса по несуществующей ссылке" {
					return
				}
				assert.Equal(t, tt.want.location, rec.Header().Get("Location"))
			}
		})
	}
}
