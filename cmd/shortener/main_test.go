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
	"testing"
)

const testUrl = "https://practicum.yandex.ru/"

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
			name: "Тест сохранения сокращенной ссылки",
			body: []byte(testUrl),
			want: want{
				statusCode:  201,
				contentType: "text/plain; charset=UTF-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(tt.body))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			requestDump, err := httputil.DumpRequest(req, true)
			if err != nil {
				panic(err)
			}
			t.Logf("Лог HTTP запроса Test_getItem %s", requestDump)

			// Assertions
			if assert.NoError(t, createItem(c)) {
				assert.Equal(t, tt.want.statusCode, rec.Code)
				assert.Equal(t, tt.want.contentType, rec.Header().Get("Content-type"))

				// проверяем body
				responseBody := rec.Body.String()
				t.Logf("Тесты пройдены %s", responseBody)

				// проверка, что в ответе url
				_, err = url.ParseRequestURI(string(responseBody))
				if err != nil {
					panic(err)
				}
				require.NoError(t, err)

				// проверка, что элемент добавлен в слайс items
				assert.Equal(t, len(items), 1)

				// получаем сокращенный url и пишем в переменную
				shortUrl = string(responseBody)
			}
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
			name: "Тест получения полной ссылки",
			url:  shortUrl,
			want: want{
				statusCode: 307,
				location:   testUrl,
			},
		},
	}

	t.Logf("Значение переменной shortUrl %s", shortUrl)
	t.Logf("Значение переменной testUrl %s", testUrl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues("557700")

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
