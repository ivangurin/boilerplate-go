package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type (
	TxFunc func(ctx context.Context, tx Executor) error

	txKey struct{}
)

const maxConns = 20

// Client предоставляет методы для работы с базой данных.
// Методы Exec, Query, QueryRow автоматически используют транзакцию,
// если она присутствует в контексте (через метод Transaction).
type Client interface {
	GetPool() *pgxpool.Pool
	Ping(ctx context.Context) error
	Transaction(ctx context.Context, fn TxFunc) error
	Close() error
	Executor
}

type Executor interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
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

func (p *client) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	tx, exists := ctx.Value(txKey{}).(pgx.Tx)
	if exists {
		return tx.Exec(ctx, sql, args...)
	}
	return p.pool.Exec(ctx, sql, args...)
}

func (p *client) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	tx, exists := ctx.Value(txKey{}).(pgx.Tx)
	if exists {
		return tx.Query(ctx, sql, args...)
	}
	return p.pool.Query(ctx, sql, args...)
}

func (p *client) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	tx, exists := ctx.Value(txKey{}).(pgx.Tx)
	if exists {
		return tx.QueryRow(ctx, sql, args...)
	}
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p *client) Transaction(ctx context.Context, fn TxFunc) error {
	conn, err := p.pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("get connection from pool: %w", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	ctx = context.WithValue(ctx, txKey{}, tx)

	err = fn(ctx, tx)
	if err != nil {
		return fmt.Errorf("execute transaction function: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
