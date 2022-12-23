package main

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

const testUrl = "test"

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
			name: "Тест успешного сохранения сокращенной ссылки",
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
			// создаём новый Recorder
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(createItem)
			// запускаем сервер
			h.ServeHTTP(w, request)
			result := w.Result()

			// проверка статус кода
			assert.Equal(t, tt.want.statusCode, result.StatusCode)

			// проверка заголовка ответа
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

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
