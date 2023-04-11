package service

import (
	"context"
	"time"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
	orderProto "github.com/RipperAcskt/innotaxiorder/pkg/proto"
)

type Service struct {
	Repo
	*OrderService
}

type Repo interface {
	CreateOrder(ctx context.Context, order model.Order) error
	GetOrders(ctx context.Context, indexes []string) ([]*model.Order, error)
	GetWaiting(ctx context.Context, taxiType string) ([]*model.Order, error)
}

func New(repo Repo, order Order, drivers []*orderProto.Driver, cfg *config.Config) *Service {
	orderService := NewOrderService(order, drivers)
	return &Service{
		repo,
		orderService,
	}
}

func (s *Service) Create(ctx context.Context, order model.Order) error {
	order.Date = time.Now().UTC().String()
	return s.CreateOrder(ctx, order)
}

func (s *Service) GetOrder(ctx context.Context, indexes []string) ([]*model.Order, error) {
	return s.GetOrders(ctx, indexes)
}

func (s *Service) Find(ctx context.Context) ([]*model.Order, error) {
	orders, err := s.GetWaiting(ctx, "econom")
	if err != nil {
		return nil, err
	}

	for i, order := range orders {
		driver := s.FindDriver(order)
		if driver == nil {
			break
		}
		orders[i].DriverID = driver.ID
		orders[i].DriverName = driver.Name
		orders[i].DriverPhone = driver.PhoneNumber
		orders[i].DriverRaiting = driver.Raiting
	}
	return orders, nil
}
