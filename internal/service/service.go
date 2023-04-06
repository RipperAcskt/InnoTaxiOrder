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
	GetOrders(ctx context.Context, indexes []string) ([]*model.Order, error)
}

func New(repo Repo, cfg *config.Config) *Service {
	return &Service{
		repo,
	}
}

func (s *Service) Create(ctx context.Context, order model.Order) error {
	return s.CreateOrder(ctx, order)
}

func (s *Service) GetOrder(ctx context.Context, indexes []string) ([]*model.Order, error) {
	return s.GetOrders(ctx, indexes)
}
