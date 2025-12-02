package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"

	"boilerplate/internal/pkg/clients/db"
)

//go:embed *.sql
var embedMigrations embed.FS

func Migrate(dbClient db.Client) error {
	goose.SetBaseFS(embedMigrations)

	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("set goose dialect: %s", err.Error())
	}

	sqlDB := sql.OpenDB(stdlib.GetPoolConnector(dbClient.GetPool()))
	defer func() {
		err = sqlDB.Close()
		if err != nil {
			panic(err)
		}
	}()

	err = goose.Up(sqlDB, ".")
	if err != nil {
		return fmt.Errorf("goose up: %s", err.Error())
	}

	return nil
}
