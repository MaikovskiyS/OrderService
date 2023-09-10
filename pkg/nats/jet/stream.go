package jet

import (
	"orderservice/internal/config"
	"orderservice/pkg/nats"

	"github.com/nats-io/nats.go/jetstream"
)

type Stream struct {
	Jet jetstream.JetStream
	Cfg *config.Nats
}

func New(cfg *config.Nats, cl *nats.Client) (*Stream, error) {
	stream, err := jetstream.New(cl.Conn, jetstream.WithClientTrace(&jetstream.ClientTrace{}))
	if err != nil {
		return nil, err
	}

	return &Stream{
		Jet: stream,
		Cfg: cfg,
	}, nil
}
