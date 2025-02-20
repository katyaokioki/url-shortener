package storage

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"sync"
)

const shortURLLength = 10

func generateShortURL() (string, error) {
	b := make([]byte, shortURLLength)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b)[:shortURLLength], nil
}

type InMemoryStorage struct {
	mu      sync.RWMutex
	data    map[string]string
	storage map[string]string
	reverse map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		storage: make(map[string]string),
		reverse: make(map[string]string),
	}
}

func (s *InMemoryStorage) Save(url string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if short, exists := s.reverse[url]; exists {
		return short, nil
	}
	short, err := generateShortURL()
	if err != nil {
		return "", err
	}
	s.storage[short] = url
	s.reverse[url] = short
	return short, nil
}

func (s *InMemoryStorage) Get(short string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if url, exists := s.storage[short]; exists {
		return url, nil
	}
	return "", errors.New("not found")
}

func (s *InMemoryStorage) Dump() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Копируем данные, чтобы не сломать основное хранилище
	copy := make(map[string]string)
	for k, v := range s.data {
		copy[k] = v
	}
	return copy
}
