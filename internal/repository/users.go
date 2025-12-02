package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"

	"boilerplate/internal/pkg/clients/db"
	"boilerplate/internal/pkg/easyscan"
)

type UserFilter struct {
	ID          []int
	Name        *string
	Email       []string
	IsAdmin     *bool
	WithDeleted *bool
	Limit       *uint64
	Offset      *uint64
	Sort        *string
}

type User struct {
	ID        int        `db:"id"`
	Name      string     `db:"name"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	IsAdmin   bool       `db:"is_admin"`
	Deleted   bool       `db:"deleted"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}

type Users struct {
	Result []*User
	Total  int
}

type UsersRepo interface {
	Create(ctx context.Context, user *User) (int, error)
	Get(ctx context.Context, id int) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, filter *UserFilter) (*Users, error)
}

type usersRepo struct {
	dbClient db.Client
}

func NewUsersRepo(dbClient db.Client) UsersRepo {
	return &usersRepo{
		dbClient: dbClient,
	}
}

func (r *usersRepo) Create(ctx context.Context, user *User) (int, error) {
	builder := sq.Insert(TableUsers).
		Columns(ColumnName, ColumnEmail, ColumnPassword, ColumnIsAdmin, ColumnCreatedAt, ColumnUpdatedAt).
		Values(user.Name, user.Email, user.Password, user.IsAdmin, "now()", "now()").
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("to sql, %s", err.Error())
	}

	var id int
	err = r.dbClient.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("execute query create user: %s", err.Error())
	}

	return id, nil
}

func (r *usersRepo) Get(ctx context.Context, id int) (*User, error) {
	builder := sq.Select("*").
		From(TableUsers).
		Where(squirrel.Eq{
			ColumnID:      id,
			ColumnDeleted: false,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("to sql, %s", err.Error())
	}

	user := &User{}
	err = easyscan.Get(ctx, r.dbClient, user, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query get user: %s", err.Error())
	}

	return user, nil
}

func (r *usersRepo) Update(ctx context.Context, user *User) error {
	builder := sq.Update(TableUsers).
		Set(ColumnName, user.Name).
		Set(ColumnEmail, user.Email).
		Set(ColumnUpdatedAt, "now()").
		Where(squirrel.Eq{
			ColumnID: user.ID,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("to sql, %s", err.Error())
	}

	_, err = r.dbClient.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execute query update user: %s", err.Error())
	}

	return nil
}

func (r *usersRepo) Delete(ctx context.Context, id int) error {
	builder := squirrel.Delete(TableUsers).
		Where(squirrel.Eq{
			ColumnID: id,
		})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("to sql, %s", err.Error())
	}

	_, err = r.dbClient.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("execute query delete user: %s", err.Error())
	}

	return nil
}

func (r *usersRepo) Search(ctx context.Context, filter *UserFilter) (*Users, error) {
	builder := sq.Select("*").
		From(TableUsers)

	if filter.ID != nil {
		builder = builder.Where(squirrel.Eq{
			ColumnID: filter.ID,
		})
	}

	if filter.Name != nil {
		builder = builder.Where(squirrel.Like{
			ColumnName: "%" + *filter.Name + "%",
		})
	}

	if filter.Email != nil {
		builder = builder.Where(squirrel.Eq{
			ColumnEmail: filter.Email,
		})
	}

	if filter.IsAdmin != nil {
		builder = builder.Where(squirrel.Eq{
			ColumnIsAdmin: *filter.IsAdmin,
		})
	}

	if filter.WithDeleted == nil {
		builder = builder.Where(squirrel.Eq{
			ColumnDeleted: false,
		})
	} else {
		builder = builder.Where(squirrel.Eq{
			ColumnDeleted: *filter.WithDeleted,
		})
	}

	if filter.Limit != nil {
		builder = builder.Limit(*filter.Limit)
	}

	if filter.Limit != nil {
		builder = builder.Limit(*filter.Limit)
	}

	if filter.Offset != nil {
		builder = builder.Offset(*filter.Offset)
	}

	if filter.Sort != nil {
		builder = builder.OrderBy(*filter.Sort)
	} else {
		builder = builder.OrderBy(ColumnID + " DESC")
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("to sql, %s", err.Error())
	}

	users := &Users{}
	err = easyscan.Select(ctx, r.dbClient, &users.Result, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query search users: %s", err.Error())
	}

	return users, nil
}
