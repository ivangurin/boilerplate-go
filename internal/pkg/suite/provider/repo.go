package suite_provider

import (
	"fmt"
	"strings"
	"sync"

	"github.com/Masterminds/squirrel"

	"boilerplate/internal/pkg/clients/db"
	"boilerplate/internal/pkg/easyscan"
	"boilerplate/internal/repository"
	"boilerplate/migrations"
)

var muRepo sync.Mutex

func (sp *Provider) GetRepo() repository.Repo {
	if sp.repo == nil {
		muRepo.Lock()

		dbClient, err := db.New(sp.ctx, sp.config.DB.GetDSN())
		if err != nil {
			panic(err)
		}

		err = migrations.Migrate(dbClient)
		if err != nil {
			panic(err)
		}

		sp.repo = repository.NewRepo(dbClient)

		sp.ClearDB()

		sp.cleanups = append(sp.cleanups,
			func() error {
				defer muRepo.Unlock()
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

	tables := []Table{}
	err = easyscan.Select(sp.Context(), sp.repo.DbClient(), &tables, sql, args...)
	if err != nil {
		panic(err)
	}

	tableNames := make([]string, len(tables))
	for i, table := range tables {
		tableNames[i] = table.Name
	}

	_, err = sp.repo.DbClient().Exec(sp.Context(), fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", strings.Join(tableNames, ", ")))
	if err != nil {
		panic(err)
	}
}
