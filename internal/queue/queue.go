package queue

import (
	"sync"

	"github.com/RipperAcskt/innotaxi/pkg/proto"
)

type Queue struct {
	Drivers    []*proto.Driver
	lastDriver map[string]*proto.Driver
	mu         sync.RWMutex
}

func New() *Queue {
	return &Queue{
		Drivers:    make([]*proto.Driver, 0),
		lastDriver: map[string]*proto.Driver{},
	}
}

func (q *Queue) Append(driver *proto.Driver) {
	last := q.getLastDriver(driver.TaxiType)
	if last != nil {
		if last.ID == driver.ID {
			return
		}
	}
	for _, d := range q.Drivers {
		if d.ID == driver.ID {
			return
		}
	}
	q.Drivers = append(q.Drivers, driver)
}

func (q *Queue) Get() *proto.Driver {
	if len(q.Drivers) == 0 {
		return nil
	}
	driver := q.Drivers[0]
	q.Drivers = q.Drivers[1:]
	q.setLastDriver(driver.TaxiType, driver)
	return driver
}

func (q *Queue) setLastDriver(taxiType string, driver *proto.Driver) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.lastDriver[taxiType] = driver
}

func (q *Queue) getLastDriver(taxiType string) *proto.Driver {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.lastDriver[taxiType]
}
