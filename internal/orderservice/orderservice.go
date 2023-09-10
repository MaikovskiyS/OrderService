package orderservice

import (
	"context"
	"orderservice/internal/orderservice/adapter/cache"
	"orderservice/internal/orderservice/adapter/storage"
	"orderservice/internal/orderservice/controller/http/handler"
	"orderservice/internal/orderservice/controller/nats-stream/msgbroker"
	"orderservice/internal/orderservice/service"
	"orderservice/pkg/nats/jet"
	"orderservice/pkg/postgres"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type Deps struct {
	Ctx      context.Context
	Logger   *logrus.Logger
	Router   *chi.Mux
	Stream   *jet.Stream
	DbClient *postgres.Client
}

// OrderService
func New(d *Deps) service.OrderService {

	store := storage.New(d.DbClient, d.Logger)
	cache := cache.New()
	svc := service.New(store, cache, d.Logger)

	consumer, err := msgbroker.New(d.Stream, svc, d.Logger, d.Ctx)
	if err != nil {
		d.Logger.Warning("cant create consumer")
		d.Ctx.Done()
		return nil
	}
	consumer.CreateOrder()

	handler := handler.New(svc)
	handler.RegisterRoutes(d.Router)

	return svc
}
