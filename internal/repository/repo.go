package repository

import (
	"github.com/Masterminds/squirrel"

	"boilerplate/internal/pkg/clients/db"
)

type Repo interface {
	DbClient() db.Client
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

func (r *repo) DbClient() db.Client {
	return r.dbClient
}

func (r *repo) Users() UsersRepo {
	if r.usersRepo == nil {
		r.usersRepo = NewUsersRepo(r.dbClient)
	}
	return r.usersRepo
}
