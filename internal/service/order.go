package service

import (
	"context"
	"fmt"

	"github.com/RipperAcskt/innotaxiorder/internal/model"
	"github.com/RipperAcskt/innotaxiorder/internal/queue"
	orderProto "github.com/RipperAcskt/innotaxiorder/pkg/proto"
)

type Order interface {
	SyncDriver(ctx context.Context, drivers []*orderProto.Driver) ([]*orderProto.Driver, error)
}

type OrderService struct {
	Order
	drivers map[string]*queue.Queue
	New     chan *orderProto.Driver
	Err     chan error
	Find    chan *model.Order
}

func NewOrderService(order Order, drivers []*orderProto.Driver) *OrderService {
	os := OrderService{
		order,
		make(map[string]*queue.Queue),
		make(chan *orderProto.Driver),
		make(chan error),
		make(chan *model.Order)}
	os.drivers[model.Econom] = queue.New()
	os.drivers[model.Comfort] = queue.New()
	os.drivers[model.Business] = queue.New()
	fmt.Println(os.drivers)
	for _, driver := range drivers {
		tmp := driver
		os.drivers[driver.TaxiType].Append(tmp)
	}
	return &os
}

func (o *OrderService) FindDriver(order *model.Order) *orderProto.Driver {
	return o.drivers[order.TaxiType].Get()
}
