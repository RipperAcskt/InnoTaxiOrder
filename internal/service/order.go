package service

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxiorder/internal/client"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
)

type Order interface {
	FindDriver(ctx context.Context, order client.OrderRequest) (*model.Order, error)
}

type OrderService struct {
	Order
	econom []client.OrderRequest
	comfort []
	business []
	New    chan client.OrderRequest
	Err    chan error
	Find   chan *model.Order
	f      bool
}

func NewOrderService(order Order) *OrderService {
	return &OrderService{order, make([]client.OrderRequest, 0), make(chan client.OrderRequest), make(chan error), make(chan *model.Order), false}
}

func (o *OrderService) Consumer() {
	for {
		select {
		case order := <-o.New:
			o.orders = append(o.orders, order)
		default:
			if len(o.orders) > 0 {
				var order client.OrderRequest
				if o.f {
					order = o.orders[len(o.orders)-1]
				}
				ord, err := o.FindDriver(context.Background(), order)
				if err != nil {
					o.f = false
					o.Err <- fmt.Errorf("find driver failed: %w", err)
				} else {
					o.f = true
					o.Find <- ord
				}
			}
		}

	}
}
