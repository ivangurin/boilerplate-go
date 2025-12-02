package repository

import (
	"github.com/Masterminds/squirrel"

	"boilerplate/internal/pkg/clients/db"
)

type Repo struct {
	dbClient  db.Client
	usersRepo UsersRepo
}

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

func New(dbClient db.Client) *Repo {
	return &Repo{
		dbClient: dbClient,
	}
}

func (r *Repo) Users() UsersRepo {
	if r.usersRepo == nil {
		r.usersRepo = NewUsersRepo(r.dbClient)
	}
	return r.usersRepo
}
