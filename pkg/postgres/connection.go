package postgres

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultMaxPoolSize  = 10
	defaultConnAttempts = 3
	defaultConnTimeout  = time.Second
)

// Postgres -.
type Client struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration
	Builder      squirrel.StatementBuilderType
	Pool         *pgxpool.Pool
}

// New -.
func New(url string) (*Client, error) {
	pg := &Client{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeout:  defaultConnTimeout,
	}
	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("postgres - NewPostgres - pgxpool.ParseConfig: %w", err)
	}
	pg.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	poolConfig.MaxConns = int32(pg.maxPoolSize)

	for pg.connAttempts > 0 {
		pg.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
		if err == nil {
			break
		}

		log.Printf("Postgres is trying to connect, attempts left: %d", pg.connAttempts)

		time.Sleep(pg.connTimeout)

		pg.connAttempts--
	}

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("postgres - NewPostgres - connAttempts == 0: %w", err)
	}

	return pg, nil
}

// Close -.
func (p *Client) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
