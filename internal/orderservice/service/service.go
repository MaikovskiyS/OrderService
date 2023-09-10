package service

import (
	"context"
	"orderservice/internal/orderservice/model"
	"time"

	"github.com/sirupsen/logrus"
)

type OrderService interface {
	GetById(ctx context.Context, id uint64) (*model.Order, error)
	Create(ctx context.Context, order *model.Order) error
}
type service struct {
	log     *logrus.Logger
	storage Storage
	cache   Cacher
}

// Constructor
func New(storage Storage, cache Cacher, l *logrus.Logger) OrderService {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		orders, err := storage.GetAll(ctx)
		if err != nil {
			l.Info("empty db")
			cancel()
			return
		}
		for _, order := range orders {
			cache.Set(order)
		}
		l.Info("cache filling is ready")
		cancel()
	}()

	return &service{
		log:     l,
		cache:   cache,
		storage: storage,
	}

}

// Get order by Id from cache, if not exist, from db and set in cache
func (s *service) GetById(ctx context.Context, id uint64) (*model.Order, error) {
	order, found := s.cache.Get(id)
	if !found {
		s.log.Infof("id: %v not found in cache", id)
		order, err := s.storage.GetById(ctx, id)
		if err != nil {
			s.log.Info("Err form db:", err)
			return &model.Order{}, err
		}
		order.OrderId = id
		s.cache.Set(order)
		return order, nil
	}

	return order, nil
}

// Create order in db and in cache
func (s *service) Create(ctx context.Context, order *model.Order) error {

	id, err := s.storage.Create(order)
	if err != nil {
		s.log.Info("Create err:", err)
		return err
	}
	order.OrderId = id
	s.cache.Set(order)
	s.log.Info("order created id:", id)
	return nil
}
