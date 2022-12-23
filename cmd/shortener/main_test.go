package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
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
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(tt.body))

			requestDump, err := httputil.DumpRequest(request, true)
			if err != nil {
				panic(err)
			}
			t.Logf("Лог запроса Test_getItem %s", requestDump)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(createItem)
			// запускаем сервер
			h.ServeHTTP(w, request)
			result := w.Result()

			t.Logf("Лог ответа Test_getItem %v", result)

			// проверка статус кода
			require.Equal(t, tt.want.statusCode, result.StatusCode)

			// проверка заголовка ответа
			require.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			responseBody, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			t.Logf("Получение ответа на запрос получения короткого url, body %s", string(responseBody))

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
		})
	}
}

func Test_getItem(t *testing.T) {
	type want struct {
		statusCode  int
		response    string
		contentType string
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
				statusCode:  200,
				contentType: "text/plain; charset=utf-8",
			},
		},
	}

	t.Logf("Значение переменной shortUrl %s", shortUrl)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, shortUrl, nil)

			requestDump, err := httputil.DumpRequest(request, false)
			if err != nil {
				panic(err)
			}

			t.Logf("Лог запроса Test_getItem %s", requestDump)

			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(createItem)
			// запускаем сервер
			h.ServeHTTP(w, request)
			result := w.Result()

			t.Logf("Лог ответа Test_getItem %v", result)

			// проверка статус кода
			require.Equal(t, tt.want.statusCode, result.StatusCode)

			//Получаем заголовок location
			location := result.Header.Get("Location")

			// проверка, что в заголовке url
			_, err = url.ParseRequestURI(location)
			if err != nil {
				panic(err)
			}

			// проверка изначально записанного url и полученного
			require.Equal(t, testUrl, location)
		})
	}
}
