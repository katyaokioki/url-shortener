package main

import (
	"log"
	"net"

	"url-shortener/api/url-shortener/api"
	"url-shortener/internal/handler"
	"url-shortener/internal/storage"

	"google.golang.org/grpc"
)

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Используем in-memory хранилище (можно заменить на Postgres)
	store := storage.NewInMemoryStorage()
	server := grpc.NewServer()
	api.RegisterURLShortenerServer(server, handler.NewGRPCServer(store))

	log.Println("gRPC Server listening on :50051")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
