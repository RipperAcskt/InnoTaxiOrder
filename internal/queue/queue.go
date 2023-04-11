package queue

import orderProto "github.com/RipperAcskt/innotaxiorder/pkg/proto"

type Queue struct {
	drivers []*orderProto.Driver
}

func New() *Queue {
	return &Queue{make([]*orderProto.Driver, 0)}
}

func (q *Queue) Append(driver *orderProto.Driver) {
	q.drivers = append(q.drivers, driver)
}

func (q *Queue) Get() *orderProto.Driver {
	if len(q.drivers) == 0 {
		return nil
	}
	driver := q.drivers[0]
	q.drivers = q.drivers[1:]
	return driver
}
