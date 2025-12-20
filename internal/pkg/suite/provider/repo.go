package suite_provider

import (
	"context"
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/gofrs/flock"
	"github.com/jackc/pgx/v5"

	"boilerplate/internal/pkg/clients/db"
	"boilerplate/internal/repository"
	"boilerplate/migrations"
)

const dbLockPath = "/db.lock"

var dbLockFile = flock.New(getLocalDirPath() + dbLockPath)

func (sp *Provider) GetRepo() repository.Repo {
	if sp.repo == nil {
		lockDB(sp.Context())

		dbClient, err := db.New(sp.ctx, sp.logger, sp.config.DB.GetDSN())
		if err != nil {
			panic(err)
		}

		err = migrations.Migrate(sp.Context(), dbClient)
		if err != nil {
			panic(err)
		}

		sp.repo = repository.NewRepo(dbClient)

		sp.ClearDB()

		sp.cleanups = append(sp.cleanups,
			func() error {
				defer unlockDB()
				// sp.ClearDB()
				sp.repo = nil
				return nil
			},
			dbClient.Close,
		)
	}

	return sp.repo
}

func (sp *Provider) ClearDB() {
	if sp.repo == nil {
		return
	}

	type Table struct {
		Name string `db:"table_name"`
	}

	builder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		Select("table_name").
		From("information_schema.tables").
		Where(
			squirrel.And{
				squirrel.Eq{"table_schema": "public"},
				squirrel.Or{
					squirrel.Eq{"table_type": "BASE TABLE"},
					squirrel.Eq{"table_type": "VIEW"},
				},
				squirrel.NotEq{"table_name": "rel_schema_versions"},
				squirrel.NotEq{"table_name": "goose_db_version"},
			},
		)

	sql, args, err := builder.ToSql()
	if err != nil {
		panic(err)
	}

	rows, err := sp.repo.Client().Query(sp.Context(), sql, args...)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	tables, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[Table])
	if err != nil {
		panic(err)
	}

	tableNames := make([]string, len(tables))
	for i, table := range tables {
		tableNames[i] = table.Name
	}

	_, err = sp.repo.Client().Exec(sp.Context(), fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(tableNames, ", ")))
	if err != nil {
		panic(err)
	}
}

func lockDB(ctx context.Context) {
	lockCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	locked, err := dbLockFile.TryLockContext(lockCtx, time.Millisecond*5)
	if err != nil {
		panic(err)
	}
	if !locked {
		panic("lockDb can't take a lock")
	}
}

func unlockDB() {
	err := dbLockFile.Unlock()
	if err != nil {
		panic(err)
	}
}

func getLocalDirPath() string {
	_, f, _, _ := runtime.Caller(0)
	return path.Dir(f)
}
