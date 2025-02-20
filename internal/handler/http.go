package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"url-shortener/internal/storage"
)

// Структура Server с публичным полем Storage
type Server struct {
	Storage storage.Storage // Публичное поле для доступа из других пакетов
}

// Конструктор для создания нового сервера
func NewServer(storage storage.Storage) *Server {
	return &Server{Storage: storage} // Инициализация публичного поля Storage
}

// Обработчик для сокращения URL
func (s *Server) HandleShorten(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}
	// Декодируем запрос в структуру
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	// Сохраняем оригинальный URL и получаем короткую ссылку
	short, err := s.Storage.Save(req.URL)
	if err != nil {
		http.Error(w, "could not save URL", http.StatusInternalServerError)
		return
	}
	// Отправляем ответ с короткой ссылкой
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"short": short})
}

// Обработчик для расширения короткой ссылки в оригинальный URL
func (s *Server) HandleExpand(w http.ResponseWriter, r *http.Request) {
	short := r.URL.Query().Get("short")
	if short == "" {
		http.Error(w, "missing short URL", http.StatusBadRequest)
		return
	}
	// Получаем оригинальный URL по короткой ссылке
	url, err := s.Storage.Get(short)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	// Отправляем ответ с оригинальным URL
	json.NewEncoder(w).Encode(map[string]string{"url": url})
}
