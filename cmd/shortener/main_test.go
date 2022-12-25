package main

import (
	"bytes"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"strings"
	"testing"
)

const testUrl = "https://practicum.yandex.ru/"
const testUrl2 = "https://www.google.com/search?q=goland+%D1%83%D1%80%D0%BE%D0%BA%D0%B8&oq=goland+%D1%83%D1%80%D0%BE%D0%BA%D0%B8&aqs=chrome..69i57j0i10i512.3638j0j15&sourceid=chrome&ie=UTF-8"

var shortUrl string

func Test_createItem(t *testing.T) {

	type want struct {
		statusCode  int
		response    string
		contentType string
	}

	tests := []struct {
		name string
		body []byte
		want want
	}{
		{
			name: "Обычная ссылка",
			body: []byte(testUrl),
			want: want{
				statusCode:  201,
				contentType: "text/plain; charset=UTF-8",
			},
		},
		{
			name: "Длинная ссылка",
			body: []byte(testUrl2),
			want: want{
				statusCode:  201,
				contentType: "text/plain; charset=UTF-8",
			},
		},
	}

	for index, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(tt.body))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			requestDump, err := httputil.DumpRequest(req, true)
			if err != nil {
				panic(err)
			}
			t.Logf("Лог HTTP запроса Test_getItem %s", requestDump)

			// Проверки
			if assert.NoError(t, createItem(c)) {
				assert.Equal(t, tt.want.statusCode, rec.Code)
				assert.Equal(t, tt.want.contentType, rec.Header().Get("Content-type"))

				// Получаем body ответа
				responseBody := rec.Body.String()
				t.Logf("Ответ сервера %s", responseBody)

				// Проверка, что в ответе url
				_, err = url.ParseRequestURI(string(responseBody))
				if err != nil {
					panic(err)
				}
				require.NoError(t, err)

				// Проверка слайса items
				assert.Equal(t, len(items), index+1)

				// Получаем сокращенный url и пишем в переменную
				shortUrl = responseBody
			}

			t.Logf("Итого элементов в слайсе items %d", len(items))

		})
	}
}

func Test_getItem(t *testing.T) {
	type want struct {
		statusCode int
		response   string
		location   string
		//contentType string
	}

	tests := []struct {
		name string
		url  string
		want want
	}{
		{
			name: "Тест получения обычной ссылки",
			url:  shortUrl,
			want: want{
				statusCode: 307,
				location:   testUrl2,
			},
		},
		//{
		//	name: "Тест получения длинной ссылки",
		//	url:  shortUrl,
		//	want: want{
		//		statusCode: 307,
		//		location:   testUrl2,
		//	},
		//},
	}

	//t.Logf("Значение переменной shortUrl %s", shortUrl)
	//t.Logf("Значение переменной testUrl2 %s", testUrl2)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")

			split := strings.Split(shortUrl, "/")
			splitLen := len(split)
			id := split[splitLen-1]

			t.Logf("Значение переменной split %s", split)
			t.Logf("Значение переменной id %s", id)

			c.SetParamValues(id)

			requestDump, err := httputil.DumpRequest(req, true)
			if err != nil {
				panic(err)
			}
			t.Logf("Лог HTTP запроса Test_getItem %s", requestDump)

			// Assertions
			if assert.NoError(t, getItem(c)) {
				assert.Equal(t, tt.want.statusCode, rec.Code)
				assert.Equal(t, tt.want.location, rec.Header().Get("Location"))
			}
		})
	}
}
