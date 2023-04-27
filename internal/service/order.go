package service

import (
	"github.com/RipperAcskt/innotaxi/pkg/proto"
	"github.com/RipperAcskt/innotaxiorder/internal/model"
	"github.com/RipperAcskt/innotaxiorder/internal/queue"
)

type OrderService struct {
	driversQueue map[model.ClassType]*queue.Queue
	Push         chan *proto.Driver
	Pop          map[model.ClassType]chan *proto.Driver
}

func newOrderService() *OrderService {
	os := OrderService{
		driversQueue: map[model.ClassType]*queue.Queue{},
		Push:         make(chan *proto.Driver),
		Pop:          map[model.ClassType]chan *proto.Driver{},
	}

	os.driversQueue[model.Econom] = queue.New()
	os.driversQueue[model.Comfort] = queue.New()
	os.driversQueue[model.Business] = queue.New()

	os.Pop[model.Econom] = make(chan *proto.Driver)
	os.Pop[model.Comfort] = make(chan *proto.Driver)
	os.Pop[model.Business] = make(chan *proto.Driver)
	return &os
}

func (o *OrderService) Append() {
	for {
		driver := <-o.Push
		taxiType := model.NewClassType(driver.TaxiType)
		o.driversQueue[taxiType].Append(driver)
	}
}

func (o *OrderService) Get() {
	go func() {
		for {
			d := o.driversQueue[model.Econom].Get()
			o.Pop[model.Econom] <- d
		}
	}()
	go func() {
		for {
			d := o.driversQueue[model.Comfort].Get()
			o.Pop[model.Comfort] <- d
		}
	}()
	go func() {
		for {
			d := o.driversQueue[model.Business].Get()
			o.Pop[model.Business] <- d
		}
	}()
}

func (o *OrderService) findDriver(order *model.Order) *proto.Driver {
	taxiType := model.NewClassType(order.TaxiType)
	return <-o.Pop[taxiType]
}
