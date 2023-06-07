package service

import (
	"context"
	"fmt"
	"time"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxiorder/config"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
)

//go:generate mockgen -destination=mocks/mock_service.go -package=mocks github.com/RipperAcskt/innotaxiorder/internal/service Repo
//go:generate mockgen -destination=mocks/mock_driver.go -package=mocks github.com/RipperAcskt/innotaxiorder/internal/service DriverService

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
	GetOrderByFilter(ctx context.Context, filters model.OrderFilters, pagginationInfo model.PagginationInfo) ([]*model.Order, error)
}

type DriverService interface {
	SyncDriver(ctx context.Context, drivers []*proto.Driver) ([]*proto.Driver, error)
	SetRating(ctx context.Context, raiting *proto.Rating) error
}

func New(repo Repo, driver DriverService, cfg *config.Config) *Service {
	orderService := newOrderService()
	service := &Service{
		driver,
		repo,
		orderService,
		make(chan error),
		cfg,
	}

	go service.Append()
	go service.Get()
	go service.TimeSync()
	return service
}

func (s *Service) Create(ctx context.Context, order model.Order) error {
	order.Date = time.Now().Format("2006-01-02 15:04:05")
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

func (s *Service) SyncDrivers(ctx context.Context, drivers []*proto.Driver) error {
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
		return nil, fmt.Errorf("you haven't any waiting orders")
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
			taxiType := model.NewClassType(order.TaxiType)
			err := s.SyncDrivers(ctx, s.driversQueue[taxiType].Drivers)
			if err != nil {
				return nil, fmt.Errorf("sync drivers failed: %w", err)
			}
			break
		}

		order.DriverID = driver.ID
		order.DriverName = driver.Name
		order.DriverPhone = driver.PhoneNumber
		order.DriverRating = float64(driver.Rating)
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

func (s *Service) GetOrdersList(ctx context.Context, filters model.OrderFilters, pagginationInfo model.PagginationInfo) ([]*model.Order, error) {
	return s.GetOrderByFilter(ctx, filters, pagginationInfo)
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

	driver := &proto.Driver{
		ID:          order.DriverID,
		Name:        order.DriverName,
		PhoneNumber: order.DriverPhone,
		Rating:      float32(order.DriverRating),
		TaxiType:    order.TaxiType,
	}

	s.Push <- driver
	return order, nil
}

func (s *Service) SetRatingService(ctx context.Context, input model.Rating, userType string, userID string) (string, error) {
	var orders []*model.Order
	var err error

	userT := model.NewUserType(userType)
	if userT == model.User {
		orders, err = s.GetOrdersByUserID(ctx, userID, model.StatusFinished.String())
		if err != nil {
			return "", fmt.Errorf("get orders failed: %w", err)
		}
	} else if userT == model.Driver {
		filters := model.OrderFilters{
			DriverID: userID,
		}
		paggination := model.PagginationInfo{
			PagginationFlag: false,
		}
		orders, err = s.GetOrderByFilter(ctx, filters, paggination)
		if err != nil {
			return "", fmt.Errorf("get orders failed: %w", err)
		}
	}

	for i, order := range orders {
		if i >= 5 {
			break
		}
		if order.ID == input.ID {
			var id string
			userT := model.NewUserType(userType)
			if userT == model.User {
				id = order.DriverID
			} else {
				id = order.UserID
			}
			rating := proto.Rating{
				Type: userType,
				ID:   id,
				Mark: float32(input.Rating),
			}
			return "", s.SetRating(ctx, &rating)
		}
	}
	return "", fmt.Errorf("order not found")
}

func (s *Service) CancelOrder(ctx context.Context, userID string) (*model.Order, error) {
	orders, err := s.GetOrders(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("get orders by id failed: %w", err)
	}
	for _, order := range orders {
		if order.UserID == userID && order.Status != model.StatusFinished && order.Status != model.StatusCanceled {
			order.Status = model.StatusCanceled
			err = s.UpdateOrder(ctx, order)
			if err != nil {
				return nil, fmt.Errorf("update order failed: %w", err)
			}
			return order, nil
		}
	}
	return nil, fmt.Errorf("you haven't any orders")
}
