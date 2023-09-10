package app

import (
	"context"
	_ "orderservice/docs"
	"orderservice/internal/config"
	"orderservice/internal/ordergenerator"
	"orderservice/internal/orderservice"
	"orderservice/internal/server"
	"orderservice/pkg/nats"
	"orderservice/pkg/nats/jet"
	"orderservice/pkg/postgres"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// Run app
func Run(cfg *config.Config) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	wg := &sync.WaitGroup{}

	//init service deps
	log := logrus.New()
	router := chi.NewMux()
	dbclient, err := postgres.New(cfg.OrderService.DbConnUrl())
	if err != nil {
		log.Warning("App-Run-PsqlClient-NewErr", err)
		return
	}
	natsclient, err := nats.New(cfg.OrderService.Nats)
	if err != nil {
		log.Warning("App-Run-NatsClient-NewErr", err)
		return
	}
	stream, err := jet.New(&cfg.OrderService.Nats, natsclient)
	if err != nil {
		log.Warning("App-Run-JetStream-NewErr:", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	log.Info("deps created")

	//swagger ui
	router.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL(cfg.SwaggerUrl)))

	deps := &orderservice.Deps{
		Ctx:      ctx,
		Logger:   log,
		DbClient: dbclient,
		Stream:   stream,
		Router:   router,
	}
	//init order service
	_ = orderservice.New(deps)

	//init order generator
	fakeorders := ordergenerator.New(stream, log)
	go func(chan os.Signal) {
		for {
			select {
			case <-ctx.Done():
				log.Println("generator gorutine stop by ctx")
				return
			case <-c:
				log.Println("generator gorutine stop by signal")
				return
			default:
				err := fakeorders.SendOrder()
				if err != nil {
					log.Warning("OrderService-FakeOrders-GenerateErr:", err)
					ctx.Done()
					return
				}
			}
			time.Sleep(time.Second * 10)
		}
	}(c)

	//init server
	server := server.New(router, cfg)
	server.Run()

	//gracefull shutdown
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		select {
		case <-c:
			log.Info("shutdown from signal")
		case <-ctx.Done():
			log.Info("shutdown from ctx")
		}
		cancel()
		log.Info("app ctx canceled")
		natsclient.Close()
		dbclient.Close()
		server.Stop()
		log.Info("deps closed")

	}(ctx)
	wg.Wait()
}
