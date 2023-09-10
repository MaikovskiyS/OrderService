package nats

import (
	"orderservice/internal/config"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	Conn   *nats.Conn
	Stream jetstream.JetStream
	JetCfg jetstream.StreamConfig
}

func New(cfg config.Nats) (*Client, error) {
	conn, err := nats.Connect(cfg.ConnUrl(), nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second))
	if err != nil {
		return &Client{}, err
	}

	stream, err := jetstream.New(conn)
	if err != nil {
		return &Client{}, err
	}
	return &Client{
		Conn:   conn,
		Stream: stream,
		JetCfg: jetstream.StreamConfig{},
	}, nil

}
func (n *Client) Close() {
	n.Conn.Close()
}
