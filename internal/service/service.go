package service

import (
	"context"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
)

type Service struct {
	Repo
	*OrderService
}

type Repo interface {
	CreateOrder(ctx context.Context, order model.Order) (string, error)
	GetOrders(ctx context.Context, indexes []string) ([]*model.Order, error)
}

func New(repo Repo, order Order, cfg *config.Config) *Service {
	return &Service{
		repo,
		NewOrderService(order),
	}
}

func (s *Service) Create(ctx context.Context, order model.Order) (string, error) {
	return s.CreateOrder(ctx, order)
}

func (s *Service) GetOrder(ctx context.Context, indexes []string) ([]*model.Order, error) {
	return s.GetOrders(ctx, indexes)
}
