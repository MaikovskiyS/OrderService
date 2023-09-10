package cache

import (
	"orderservice/internal/orderservice/model"
	"orderservice/internal/orderservice/service"
	"sync"
	"time"
)

type cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	store             map[uint64]Item
}
type Item struct {
	Order      *model.Order
	Created    time.Time
	Expiration int64
}

// Constructor
func New() service.Cacher {
	return &cache{
		store:             make(map[uint64]Item),
		defaultExpiration: 30 * time.Second,
		cleanupInterval:   30 * time.Second,
	}
}

// Set Order model to cache
func (c *cache) Set(order *model.Order) {
	expiration := time.Now().Add(c.defaultExpiration).UnixNano()
	c.RLock()
	defer c.RUnlock()
	c.store[order.OrderId] = Item{
		Order:      order,
		Expiration: expiration,
		Created:    time.Now(),
	}
}

// Get model Order from cache
func (c *cache) Get(id uint64) (order *model.Order, found bool) {
	c.RLock()
	defer c.RUnlock()
	item, ok := c.store[id]
	if !ok {
		return order, false
	}
	if time.Now().UnixNano() > item.Expiration {
		c.Delete(id)
		return order, false
	}
	return item.Order, true

}

// Delete Order model from cache
func (c *cache) Delete(id uint64) {
	delete(c.store, id)
}
