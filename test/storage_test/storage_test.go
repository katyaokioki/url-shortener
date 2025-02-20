package storage_test

import (
	"testing"
	"url-shortener/internal/storage"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryStorage_SaveAndGet(t *testing.T) {
	// Создаем новое хранилище
	store := storage.NewInMemoryStorage()

	// Ожидаем, что для первого URL будет сгенерирована короткая ссылка
	originalURL := "https://example.com"
	shortURL, err := store.Save(originalURL)
	assert.NoError(t, err)
	assert.Len(t, shortURL, 10, "Short URL should have 10 characters")

	// Теперь получаем оригинальный URL по короткому
	result, err := store.Get(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, originalURL, result, "Original URL should match the stored URL")

	// Проверим, что попытка получить несуществующую ссылку вернет ошибку
	_, err = store.Get("nonexistent")
	assert.Error(t, err, "Expected error for non-existent short URL")
	t.Logf("Storage dump: %+v", store.Dump())
}
