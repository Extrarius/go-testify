package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMainHandlerWhenOk проверяет, что запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое
func TestMainHandlerWhenOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()

	require.Equal(t, responseRecorder.Code, http.StatusOK)
	assert.NotEmpty(t, body)
}

// TestMainHandlerWhenMissingCount проверяет, что город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerWhenMissingCount(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=7&city=tula", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	body := responseRecorder.Body.String()

	require.Equal(t, responseRecorder.Code, http.StatusBadRequest)

	expected := "wrong city value"
	assert.Equal(t, expected, body)
}

// если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе. Если сервис вернул статус код, отличный от 200, нужно остановить выполнение теста. В сервисе будет только один город moscow, в котором будет всего 4 кафе
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, responseRecorder.Code, http.StatusOK)

	body := responseRecorder.Body.String()
	list := strings.Split(body, ",")

	assert.Len(t, list, totalCount)
}
