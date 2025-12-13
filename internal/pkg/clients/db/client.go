package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxConns = 20

type Client interface {
	GetPool() *pgxpool.Pool
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
	Ping(ctx context.Context) error
	Close() error
}

type client struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (Client, error) {
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("db parse config: %w", err)
	}

	config.MaxConns = maxConns
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute

	p, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %w", err)
	}

	err = p.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return &client{
		pool: p,
	}, nil
}

func (p *client) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return p.pool.Exec(ctx, sql, args...)
}

func (p *client) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return p.pool.Query(ctx, sql, args...)
}

func (p *client) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p *client) BeginFunc(ctx context.Context, f func(pgx.Tx) error) error {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("get connection from pool: %v", err)
	}
	defer conn.Release()

	err = pgx.BeginFunc(ctx, conn, f)
	if err != nil {
		return fmt.Errorf("make transaction: %v", err)
	}

	return nil
}

func (p *client) GetPool() *pgxpool.Pool {
	return p.pool
}

func (p *client) Ping(ctx context.Context) error {
	return p.pool.Ping(ctx)
}

func (p *client) Close() error {
	p.pool.Close()
	return nil
}
