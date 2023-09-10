package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SwaggerUrl   string
	OrderService OrderService
}
type OrderService struct {
	httpPort  string
	dbConnUrl string
	Nats      Nats
}
type Nats struct {
	connUrl       string
	streamName    string
	streamSubject []string
}

// Constructor
func New() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return &Config{}, errors.New("no .env file found")
	}
	natsStreamSubject := make([]string, 1)
	port := os.Getenv("HTTP_PORT")
	pgurl := os.Getenv("DBCONN_URL")
	natsurl := os.Getenv("NATS_URL")
	natsStreamName := os.Getenv("STREAM_NAME")
	natsStreamSubject[0] = os.Getenv("STREAM_SUBJECT")
	swaggerurl := os.Getenv("SWAGGER_URL")
	return &Config{
		SwaggerUrl: swaggerurl,
		OrderService: OrderService{
			httpPort:  ":" + port,
			dbConnUrl: pgurl,
			Nats: Nats{
				connUrl:       natsurl,
				streamName:    natsStreamName,
				streamSubject: natsStreamSubject,
			},
		},
	}, nil
}

// Get server port
func (o *OrderService) Port() string {
	return o.httpPort
}

// Get db connection url
func (o *OrderService) DbConnUrl() string {
	return o.dbConnUrl
}

// Get nats connection url
func (n *Nats) ConnUrl() string {
	return n.connUrl
}

// Get jetstream name
func (n *Nats) StreamName() string {
	return n.streamName
}

// Get jetstream subject
func (n *Nats) StreamSubject() []string {
	return n.streamSubject
}
