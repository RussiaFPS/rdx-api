package main

import (
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostHandle(t *testing.T) { //Проверка post запроса
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/post/dfg/123", nil) // Добавление данных
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Successful", w.Body.String())
}

func TestUpdateNameHandle(t *testing.T) { // Проверка patch запроса на обновление имени
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PATCH", "/update/name/t1/test1", nil) // Правильный запрос
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Successful", w.Body.String())

	w = httptest.NewRecorder()

	req, _ = http.NewRequest("PATCH", "/update/name/hhrttr/test1", nil) // Неправильный запрос
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "Failed", w.Body.String())
}

func TestUpdatePassHandle(t *testing.T) { // Проверка update запроса на обновление пароля
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("PATCH", "/update/password/test1/321", nil) // Правильный запрос
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Successful", w.Body.String())

	w = httptest.NewRecorder()

	req, _ = http.NewRequest("PATCH", "/update/password/hhvvci/321", nil) // Неправильный запрос
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "Failed", w.Body.String())
}

func TestDeleteHandle(t *testing.T) { //Проверка запроса delete данных
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("DELETE", "/delete/9", nil) //Правильный запрос
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "Successful", w.Body.String())

	w = httptest.NewRecorder()

	req, _ = http.NewRequest("DELETE", "/delete/99", nil) //Неправильный запрос
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "Failed", w.Body.String())
}

func TestGetHandle(t *testing.T) { //Проверка запроса get
	router := setupRouter()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/get/1", nil) //Правильный запрос get
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"createAt\":\"2022-09-20T00:00:00Z\",\"id\":1,\"password\":\"$2a$10$JOhAj6VlaB4KA5eHiqadJOGYJC0lVYCz/hkkBcHFPuj7H9OGos7cm\",\"updateAt\":\"2022-09-22T00:00:00Z\",\"username\":\"test1\"}", w.Body.String())

	w = httptest.NewRecorder()

	req, _ = http.NewRequest("GET", "/get/99", nil) //Неправильный запрос
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, "Failed", w.Body.String())
}
