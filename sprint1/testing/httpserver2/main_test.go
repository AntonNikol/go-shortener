//package main
//
//import (
//	"github.com/stretchr/testify/assert"
//	"net/http"
//	"net/http/httptest"
//	"strings"
//	"testing"
//)
//
//func TestUserViewHandler(t *testing.T) {
//	//type args struct {
//	//	users map[string]User
//	//}
//	//type users map[string]User
//
//	users := make(map[string]User)
//
//	type want struct {
//		code        int
//		response    string
//		contentType string
//		users       map[string]User
//		url         string
//	}
//	u1 := User{
//		ID:        "u1",
//		FirstName: "Misha",
//		LastName:  "Popov",
//	}
//	u2 := User{
//		ID:        "u2",
//		FirstName: "Sasha",
//		LastName:  "Popov",
//	}
//	users["u1"] = u1
//	users["u2"] = u2
//
//	tests := []struct {
//		name  string
//		users map[string]User
//		want  want
//	}{
//		{
//			name: "Тест отсутствия query параметра user_id",
//			want: want{
//				code:  400,
//				users: users,
//				url:   "/users",
//			},
//		},
//		{
//			name: "Тест user не найден",
//			want: want{
//				code:  404,
//				users: users,
//				url:   "/users?user_id=u4",
//			},
//		},
//		{
//			name: "Тест user найден",
//			want: want{
//				code:        200,
//				users:       users,
//				url:         "/users?user_id=u2",
//				contentType: "application/json",
//			},
//		},
//	}
//	for _, tt := range tests {
//		// Запускаем каждый тест
//		t.Run(tt.name, func(t *testing.T) {
//			//if got := UserViewHandler(tt.want.users); !reflect.DeepEqual(got, tt.want) {
//			//	t.Errorf("UserViewHandler() = %v, want %v", got, tt.want)
//			//}
//			request := httptest.NewRequest(http.MethodGet, tt.want.url, nil)
//
//			// создаём новый Recorder
//			w := httptest.NewRecorder()
//			// определяем хендлер
//			h := UserViewHandler(users)
//			// запускаем сервер
//			h.ServeHTTP(w, request)
//
//			res := w.Result()
//
//			// Сверяем результат
//			//assert.Equal(t, res.StatusCode, tt.want.code)
//			if res.StatusCode != tt.want.code {
//				t.Errorf("Expected status code %d, got %d", tt.want.code, res.StatusCode)
//			}
//
//			if strings.Contains(tt.name, "Тест user найден") {
//				t.Logf("Проверка что ответ JSON если юзер найден")
//				assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
//			}
//		})
//	}
//}

//v2 с проверкой декодированного ответа

package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserViewHandler(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
		user        User
	}
	tests := []struct {
		name    string
		request string
		users   map[string]User
		want    want
	}{
		{
			name: "simple test #1",
			users: map[string]User{
				"id1": {
					ID:        "id1",
					FirstName: "Misha",
					LastName:  "Popov",
				},
			},
			want: want{
				contentType: "application/json",
				statusCode:  200,
				user: User{ID: "id1",
					FirstName: "Misha",
					LastName:  "Popov",
				},
			},
			request: "/users?user_id=id1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(UserViewHandler(tt.users))
			h(w, request)
			result := w.Result()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
			userResult, err := ioutil.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
			var user User
			err = json.Unmarshal(userResult, &user)
			require.NoError(t, err)
			assert.Equal(t, tt.want.user, user)
		})
	}
}
