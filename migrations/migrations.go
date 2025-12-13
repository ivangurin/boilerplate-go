package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"boilerplate/internal/pkg/clients/db"
)

//go:embed *.sql
var migrations embed.FS

func Migrate(ctx context.Context, dbClient db.Client) error {
	goose.SetBaseFS(migrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	sqlDB := sql.OpenDB(stdlib.GetPoolConnector(dbClient.GetPool()))
	defer func() {
		err = sqlDB.Close()
		if err != nil {
			panic(err)
		}
	}()

	err = goose.UpContext(ctx, sqlDB, ".")
	if err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	return nil
}
