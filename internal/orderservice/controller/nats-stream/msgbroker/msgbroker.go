package msgbroker

import (
	"context"
	"encoding/json"
	"orderservice/internal/orderservice/service"
	"orderservice/pkg/nats/jet"
	"time"

	"github.com/nats-io/nats.go/jetstream"
	"github.com/sirupsen/logrus"
)

type Messager interface {
	CreateOrder()
}
type broker struct {
	ctx      context.Context
	log      *logrus.Logger
	consumer jetstream.Consumer
	svc      service.OrderService
}

// Constructor
func New(stream *jet.Stream, svc service.OrderService, l *logrus.Logger, ctx context.Context) (Messager, error) {

	c, err := stream.Jet.CreateOrUpdateConsumer(ctx, stream.Cfg.StreamName(), jetstream.ConsumerConfig{
		AckPolicy:     jetstream.AckAllPolicy,
		DeliverPolicy: jetstream.DeliverPolicy(2),
	})
	if err != nil {
		l.Info("MsgBroker-Constructor-CreateOrUpdateConsumer Err:", err)
		return nil, err
	}

	return &broker{
		ctx:      ctx,
		log:      l,
		svc:      svc,
		consumer: c,
	}, nil
}

// Get Msg from jetstream
func (b *broker) CreateOrder() {

	_, err := b.consumer.Consume(func(msg jetstream.Msg) {
		ctx, cancel := context.WithTimeout(b.ctx, 3*time.Second)
		defer cancel()
		orderDto := &OrderDTO{}

		err := json.Unmarshal((msg.Data()), &orderDto)
		if err != nil {
			b.log.Info("MsgBroker-Consume-UnmarshallErr:", err)
			msg.Ack()
			return
		}
		err = orderDto.Validate()
		if err != nil {
			b.log.Info("MsgBroker-Consume-DTOVailadation Err:", err)
			msg.Ack()

		}
		order := orderDto.toModel()

		b.svc.Create(ctx, order)
		msg.Ack()
	})
	if err != nil {
		b.log.Info("MsgBroker-ConsumeErr:", err)
		return
	}
}
