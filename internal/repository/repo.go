package repository

import (
	"context"

	"github.com/Masterminds/squirrel"

	"boilerplate/internal/pkg/clients/db"
)

type Repo interface {
	Client() db.Client
	Transaction(ctx context.Context, fn db.TxFunc) error
	Users() UsersRepo
}

type repo struct {
	dbClient  db.Client
	usersRepo UsersRepo
}

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func NewRepo(dbClient db.Client) Repo {
	return &repo{
		dbClient: dbClient,
	}
}

func (r *repo) Client() db.Client {
	return r.dbClient
}

func (r *repo) Transaction(ctx context.Context, fn db.TxFunc) error {
	return r.dbClient.Transaction(ctx, fn)
}

func (r *repo) Users() UsersRepo {
	if r.usersRepo == nil {
		r.usersRepo = NewUsersRepo(r.dbClient)
	}
	return r.usersRepo
}
