package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

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
	p, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create db pool: %s", err.Error())
	}
	err = p.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("ping db: %s", err.Error())
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
	fmt.Println("11")
	p.pool.Close()
	fmt.Println("22")
	return nil
}
