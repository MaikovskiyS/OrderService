package service

import (
	"context"
	"orderservice/internal/orderservice/model"
)

type Storage interface {
	Create(order *model.Order) (uint64, error)
	GetAll(ctx context.Context) ([]*model.Order, error)
	GetById(ctx context.Context, id uint64) (*model.Order, error)
}
type Cacher interface {
	Set(order *model.Order)
	Get(id uint64) (order *model.Order, found bool)
	Delete(id uint64)
}
