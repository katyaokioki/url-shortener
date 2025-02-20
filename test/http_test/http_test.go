package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"url-shortener/internal/handler"
	"url-shortener/internal/storage"

	"github.com/stretchr/testify/assert"
)

func TestHandleShorten(t *testing.T) {
	// Инициализируем хранилище
	store := storage.NewInMemoryStorage()
	server := handler.NewServer(store)

	// Запускаем сервер в горутине
	go func() {
		http.HandleFunc("/v1/shorten", server.HandleShorten)
		http.ListenAndServe(":8080", nil)
	}()

	// Создаем данные для теста
	url := "https://example.com"
	data := map[string]string{"url": url}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal json: %v", err)
	}

	// Отправляем POST запрос
	resp, err := http.Post("http://localhost:8080/v1/shorten", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что сервер вернул статус 201
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Декодируем ответ в map
	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Проверяем, что ответ содержит короткую ссылку
	assert.NotEmpty(t, response["short"])
	assert.Len(t, response["short"], 10) // Проверяем, что ссылка имеет длину 10
}

func TestHandleExpand(t *testing.T) {
	// Инициализируем хранилище
	store := storage.NewInMemoryStorage()
	server := handler.NewServer(store)

	// Запускаем сервер в горутине
	go func() {
		http.HandleFunc("/v1/expand", server.HandleExpand)
		http.ListenAndServe(":8080", nil)
	}()

	// Сначала создаем сокращенный URL
	url := "https://example.com"
	short, err := store.Save(url)
	if err != nil {
		t.Fatalf("Failed to shorten URL: %v", err)
	}

	// Строим URL для расширения
	resp, err := http.Get("http://localhost:8080/v1/expand?short=" + short)
	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что сервер вернул статус 200
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Декодируем ответ в map
	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Проверяем, что оригинальный URL соответствует
	assert.Equal(t, url, response["url"])
}
