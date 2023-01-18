package main

import (
	"bytes"
	"encoding/json"
	"github.com/AntonNikol/go-shortener/internal/app/repositories"
	"github.com/AntonNikol/go-shortener/internal/app/repositories/inmemory"
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

var h = handlers.New("http://localhost:8080", repositories.Repository(inmemory.New()))

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.body)))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			// Проверки
			if assert.NoError(t, h.CreateItem(c)) {
				require.Equal(t, tt.want.statusCode, rec.Code)
				require.Equal(t, tt.want.contentType, rec.Header().Get("Content-type"))

				// Получаем body ответа
				responseBody := rec.Body.String()
				//t.Logf("Ответ сервера %s", responseBody)

				// Проверка, что в ответе url
				_, err := url.ParseRequestURI(responseBody)
				require.NoError(t, err)
			}
		})
	}
}

// Тест сокращения ссылки с JsonBody
func Test_createItemJSON(t *testing.T) {

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
			body: `{"url": "https://practicum.yandex.ru/"}`,
			want: want{
				statusCode:  201,
				contentType: "application/json; charset=UTF-8",
			},
		},
		{
			name: "Длинная ссылка",
			body: `{"url": "https://www.google.com/search?q=goland+%D1%83%D1%80%D0%BE%D0%BA%D0%B8&oq=goland+%D1%83%D1%80%D0%BE%D0%BA%D0%B8&aqs=chrome..69i57j0i10i512.3638j0j15&sourceid=chrome&ie=UTF-8"}`,
			want: want{
				statusCode:  201,
				contentType: "application/json; charset=UTF-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(tt.body)))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Request().Header.Set("Content-Type", "application/json")

			// Проверки
			if assert.NoError(t, h.CreateItemJSON(c)) {
				require.Equal(t, tt.want.statusCode, rec.Code)
				require.Equal(t, tt.want.contentType, rec.Header().Get("Content-type"))

				// Получаем body ответа
				responseBody := rec.Body.String()
				//t.Logf("Ответ сервера %s", responseBody)

				// проверка что это json, декодируем в мапу
				var response map[string]string
				err := json.Unmarshal([]byte(responseBody), &response)
				require.NoError(t, err)

				// проверка что в мапе есть result
				value, exist := response["result"]
				require.True(t, exist)

				// Проверка, что это url
				_, err = url.ParseRequestURI(value)
				require.NoError(t, err)
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

	h := handlers.New("http://localhost:8080", repositories.Repository(inmemory.New()))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(tt.want.location)))
			rec := httptest.NewRecorder()
			ctx := e.NewContext(req, rec)

			h.CreateItem(ctx)
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
			if assert.NoError(t, h.GetItem(c)) {
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
