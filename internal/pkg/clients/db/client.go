package db

import (
	"boilerplate/internal/pkg/logger"
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
	logger logger.Logger
	pool   *pgxpool.Pool
}

func New(ctx context.Context, logger logger.Logger, dsn string) (Client, error) {
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
		logger: logger,
		pool:   p,
	}, nil
}

func (c *client) GetPool() *pgxpool.Pool {
	return c.pool
}

func (c *client) Ping(ctx context.Context) error {
	return c.pool.Ping(ctx)
}

func (c *client) Close() error {
	c.pool.Close()
	return nil
}

func (c *client) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	now := time.Now().UTC()
	c.logger.DebugKV(ctx, "start execution", "started at", now.Format(time.RFC3339Nano), "sql", sql, "args", fmt.Sprintf("%+v", args))
	defer c.logger.DebugKV(ctx, "end execution", "ended at", time.Now().UTC().Format(time.RFC3339Nano), "duration", fmt.Sprint(time.Since(now)))

	tx, exists := ctx.Value(txKey{}).(pgx.Tx)
	if exists {
		return tx.Exec(ctx, sql, args...)
	}
	return c.pool.Exec(ctx, sql, args...)
}

func (c *client) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	now := time.Now().UTC()
	c.logger.DebugKV(ctx, "start query", "started at", now.Format(time.RFC3339Nano), "sql", sql, "args", fmt.Sprintf("%+v", args))
	defer c.logger.DebugKV(ctx, "end query", "ended at", time.Now().UTC().Format(time.RFC3339Nano), "duration", fmt.Sprint(time.Since(now)))

	tx, exists := ctx.Value(txKey{}).(pgx.Tx)
	if exists {
		return tx.Query(ctx, sql, args...)
	}
	return c.pool.Query(ctx, sql, args...)
}

func (c *client) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	now := time.Now().UTC()
	c.logger.DebugKV(ctx, "start query row", "started at", now.Format(time.RFC3339Nano), "sql", sql, "args", fmt.Sprintf("%+v", args))
	defer c.logger.DebugKV(ctx, "end query row", "ended at", time.Now().UTC().Format(time.RFC3339Nano), "duration", fmt.Sprint(time.Since(now)))

	tx, exists := ctx.Value(txKey{}).(pgx.Tx)
	if exists {
		return tx.QueryRow(ctx, sql, args...)
	}
	return c.pool.QueryRow(ctx, sql, args...)
}

func (c *client) Transaction(ctx context.Context, fn TxFunc) error {
	now := time.Now().UTC()
	c.logger.DebugKV(ctx, "start transaction", "started at", now.Format(time.RFC3339Nano))
	defer c.logger.DebugKV(ctx, "end transaction", "ended at", time.Now().UTC().Format(time.RFC3339Nano), "duration", fmt.Sprint(time.Since(now)))

	conn, err := c.pool.Acquire(ctx)
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
