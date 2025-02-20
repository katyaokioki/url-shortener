package grpc_test

import (
	"context"
	"log"
	"net"
	"testing"
	"time"
	"url-shortener/api/url-shortener/api"
	"url-shortener/internal/handler"
	"url-shortener/internal/storage"

	"google.golang.org/grpc"
)

func startGRPCServer() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	store := storage.NewInMemoryStorage()
	server := grpc.NewServer()
	api.RegisterURLShortenerServer(server, handler.NewGRPCServer(store))

	log.Println("gRPC Server listening on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func TestShortenURL(t *testing.T) {
	// Запускаем сервер в горутине
	go startGRPCServer()
	time.Sleep(2 * time.Second) // Ждём, чтобы сервер запустился

	// Создаём gRPC-клиент
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := api.NewURLShortenerClient(conn)

	// Отправляем запрос на сокращение URL
	req := &api.ShortenRequest{OriginalUrl: "https://example.com"}
	resp, err := client.ShortenURL(context.Background(), req)
	if err != nil {
		t.Fatalf("ShortenURL failed: %v", err)
	}

	// Проверяем, что ссылка действительно сократилась
	if len(resp.ShortUrl) != 10 {
		t.Errorf("Expected short URL length 10, got %d", len(resp.ShortUrl))
	}

	t.Logf("Shortened URL: %s", resp.ShortUrl)
}
