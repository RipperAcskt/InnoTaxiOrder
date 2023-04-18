package service

import (
	"context"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
	orderProto "github.com/RipperAcskt/innotaxiorder/pkg/proto"
)

var (
	ErrNotFoud = fmt.Errorf("not found")
)

type Service struct {
	DriverService
	Repo
	*OrderService
	Err chan error
	cfg *config.Config
}

type Repo interface {
	CreateOrder(ctx context.Context, order model.Order) error
	GetOrders(ctx context.Context, indexes []string) ([]*model.Order, error)
	GetStatus(ctx context.Context, taxiType, status string) ([]*model.Order, error)
	UpdateOrder(ctx context.Context, order *model.Order) error
	GetOrdersByUserID(ctx context.Context, index string, status string) ([]*model.Order, error)
}

type DriverService interface {
	SyncDriver(ctx context.Context, drivers []*orderProto.Driver) ([]*orderProto.Driver, error)
}

func New(repo Repo, driver DriverService, cfg *config.Config) *Service {
	orderService := newOrderService()
	return &Service{
		driver,
		repo,
		orderService,
		make(chan error),
		cfg,
	}
}

func (s *Service) Create(ctx context.Context, order model.Order) error {
	order.Date = time.Now().UTC().String()
	return s.CreateOrder(ctx, order)
}

func (s *Service) GetOrder(ctx context.Context, indexes []string) ([]*model.Order, error) {
	return s.GetOrders(ctx, indexes)
}

func (s *Service) GetOrderByID(ctx context.Context, index string, status string) ([]*model.Order, error) {
	return s.GetOrdersByUserID(ctx, index, status)
}

func (s *Service) TimeSync() {
	t := time.NewTicker(time.Second * time.Duration(s.cfg.SYNC_TIME))
	for range t.C {

		err := s.SyncDrivers(context.Background(), s.driversQueue[model.Econom].Drivers)
		if err != nil {
			s.Err <- fmt.Errorf("sync drivers econom failed: %w", err)
		}

		err = s.SyncDrivers(context.Background(), s.driversQueue[model.Comfort].Drivers)
		if err != nil {
			s.Err <- fmt.Errorf("sync drivers comfort failed: %w", err)
		}

		err = s.SyncDrivers(context.Background(), s.driversQueue[model.Business].Drivers)
		if err != nil {
			s.Err <- fmt.Errorf("sync drivers buisiness failed: %w", err)
		}
	}

}

func (s *Service) SyncDrivers(ctx context.Context, drivers []*orderProto.Driver) error {
	drivers, err := s.SyncDriver(ctx, drivers)
	if err != nil {
		return fmt.Errorf("can't sync drivers: %w", err)
	}
	for _, driver := range drivers {
		tmp := driver
		s.Push <- tmp
	}
	return nil
}

func (s *Service) Find(ctx context.Context, userID string) (*model.Order, error) {
	orders, err := s.GetOrdersByUserID(ctx, userID, model.StatusWaiting.String())
	if err != nil {
		return nil, fmt.Errorf("get orders by id failed: %w", err)
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("you haven't any orders")
	}
	order := orders[0]

	var indexes []string
	indexes = append(indexes, order.ID)
	orders, err = s.GetOrders(ctx, indexes)
	if err != nil {
		return nil, fmt.Errorf("get orders failed: %w", err)
	}

	for _, order := range orders {
		if order.UserID == userID && order.Status == model.StatusFound {
			return order, nil
		}
	}

	orders, err = s.GetStatus(ctx, order.TaxiType, model.StatusWaiting.String())
	if err != nil {
		return nil, fmt.Errorf("get status failed: %w", err)
	}

	var foundOrder *model.Order
	for _, order := range orders {
		driver := s.findDriver(order)
		if driver == nil {
			s.SyncDrivers(ctx, s.driversQueue[order.TaxiType].Drivers)
			break
		}

		order.DriverID = driver.ID
		order.DriverName = driver.Name
		order.DriverPhone = driver.PhoneNumber
		order.DriverRaiting = float64(driver.Raiting)
		order.Status = model.StatusFound
		err = s.UpdateOrder(ctx, order)
		if err != nil {
			return nil, fmt.Errorf("update order failed: %w", err)
		}
		if order.UserID == userID {
			tmp := order
			foundOrder = tmp
		}
	}

	if foundOrder != nil {
		foundOrder.Status = model.StatusInProgress
		err = s.UpdateOrder(ctx, foundOrder)
		if err != nil {
			return nil, fmt.Errorf("update order failed: %w", err)
		}
		return foundOrder, nil
	}
	return nil, ErrNotFoud
}

func (s *Service) CompleteOrder(ctx context.Context, userID string) (*model.Order, error) {
	orders, err := s.GetOrdersByUserID(ctx, userID, model.StatusInProgress.String())
	if err != nil {
		return nil, fmt.Errorf("get orders by id failed: %w", err)
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("you haven't any orders")
	}
	order := orders[0]

	order.Status = model.StatusFinished
	err = s.UpdateOrder(ctx, order)
	if err != nil {
		return nil, fmt.Errorf("update order failed: %w", err)
	}

	driver := &orderProto.Driver{
		ID:          order.DriverID,
		Name:        order.DriverName,
		PhoneNumber: order.DriverPhone,
		Raiting:     float32(order.DriverRaiting),
		TaxiType:    order.TaxiType,
	}

	s.Push <- driver
	return order, nil
}
