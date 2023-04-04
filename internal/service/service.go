package service

import (
	"context"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
)

type Service struct {
	Repo
}

type Repo interface {
	CreateOrder(ctx context.Context, order model.Order) error
}

func New(repo Repo, cfg *config.Config) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) GetProfile(ctx context.Context, order model.Order) error {
	return s.CreateOrder(ctx, order)
}
