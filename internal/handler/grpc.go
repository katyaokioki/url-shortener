package handler

import (
	"context"
	"url-shortener/api/url-shortener/api"
	"url-shortener/internal/storage"
)

type GRPCServer struct {
	api.UnimplementedURLShortenerServer
	storage storage.Storage
}

func NewGRPCServer(storage storage.Storage) *GRPCServer {
	return &GRPCServer{storage: storage}
}

func (s *GRPCServer) ShortenURL(ctx context.Context, req *api.ShortenRequest) (*api.ShortenResponse, error) {
	short, err := s.storage.Save(req.OriginalUrl)
	if err != nil {
		return nil, err
	}
	return &api.ShortenResponse{ShortUrl: short}, nil
}

func (s *GRPCServer) ExpandURL(ctx context.Context, req *api.ExpandRequest) (*api.ExpandResponse, error) {
	url, err := s.storage.Get(req.ShortUrl)
	if err != nil {
		return nil, err
	}
	return &api.ExpandResponse{OriginalUrl: url}, nil
}
