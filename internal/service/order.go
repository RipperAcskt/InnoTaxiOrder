package service

import (
	"github.com/RipperAcskt/innotaxiorder/internal/model"
	"github.com/RipperAcskt/innotaxiorder/internal/queue"
	orderProto "github.com/RipperAcskt/innotaxiorder/pkg/proto"
)

type OrderService struct {
	driversQueue map[string]*queue.Queue
	Push         chan *orderProto.Driver
	Pop          map[string]chan *orderProto.Driver
}

func newOrderService() *OrderService {
	os := OrderService{
		driversQueue: map[string]*queue.Queue{},
		Push:         make(chan *orderProto.Driver),
		Pop:          map[string]chan *orderProto.Driver{},
	}

	os.driversQueue[model.Econom] = queue.New()
	os.driversQueue[model.Comfort] = queue.New()
	os.driversQueue[model.Business] = queue.New()

	os.Pop[model.Econom] = make(chan *orderProto.Driver)
	os.Pop[model.Comfort] = make(chan *orderProto.Driver)
	os.Pop[model.Business] = make(chan *orderProto.Driver)
	return &os
}

func (o *OrderService) Append() {
	for {
		driver := <-o.Push
		o.driversQueue[driver.TaxiType].Append(driver)
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

func (o *OrderService) findDriver(order *model.Order) *orderProto.Driver {
	return <-o.Pop[order.TaxiType]
}
