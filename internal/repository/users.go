package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"boilerplate/internal/pkg/clients/db"
)

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

type UserFilter struct {
	IDs         []int
	Name        *string
	Emails      []string
	IsAdmin     *bool
	WithDeleted *bool
	Limit       *int
	Offset      *int
	Sort        *string
}

type Users struct {
	Result []*User
	Total  int
}

type UsersRepo interface {
	Create(ctx context.Context, user *User) error
	Get(ctx context.Context, id int) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int) error
	Search(ctx context.Context, filter *UserFilter) (*Users, error)
}

type usersRepo struct {
	client db.Client
}

func NewUsersRepo(client db.Client) UsersRepo {
	return &usersRepo{
		client: client,
	}
}

func (r *usersRepo) Create(ctx context.Context, user *User) error {
	builder := sq.Insert(TableUsers).
		Columns(ColumnName, ColumnEmail, ColumnPassword, ColumnIsAdmin, ColumnCreatedAt, ColumnUpdatedAt).
		Values(user.Name, user.Email, user.Password, user.IsAdmin, squirrel.Expr("now()"), squirrel.Expr("now()")).
		Suffix("RETURNING *")

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("to sql: %w", err)
	}

	rows, err := r.client.Query(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("execute query create user: %w", err)
	}

	defer rows.Close()

	createdUser, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return fmt.Errorf("collect user: %w", err)
	}

	*user = *createdUser

	return nil
}

func (r *usersRepo) Get(ctx context.Context, id int) (*User, error) {
	builder := sq.Select("*").
		From(TableUsers).
		Where(squirrel.Eq{
			ColumnID: id,
		})

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("to sql: %w", err)
	}

	rows, err := r.client.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query get user: %w", err)
	}
	defer rows.Close()

	user, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, fmt.Errorf("collect user: %w", err)
	}

	return user, nil
}

func (r *usersRepo) Update(ctx context.Context, user *User) error {
	builder := sq.Update(TableUsers).
		Set(ColumnName, user.Name).
		Set(ColumnEmail, user.Email).
		Set(ColumnPassword, user.Password).
		Set(ColumnUpdatedAt, squirrel.Expr("now()")).
		Where(squirrel.Eq{
			ColumnID: user.ID,
		})

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("to sql: %w", err)
	}

	_, err = r.client.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("execute query update user: %w", err)
	}

	return nil
}

func (r *usersRepo) Delete(ctx context.Context, id int) error {
	builder := sq.Update(TableUsers).
		Set(ColumnDeleted, true).
		Set(ColumnDeletedAt, squirrel.Expr("now()")).
		Where(squirrel.Eq{
			ColumnID: id,
		})

	sql, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("to sql: %w", err)
	}

	_, err = r.client.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("execute query delete user: %w", err)
	}

	return nil
}

func (r *usersRepo) Search(ctx context.Context, filter *UserFilter) (*Users, error) {
	builder := sq.Select("*").
		From(TableUsers)

	if filter.IDs != nil {
		builder = builder.Where(squirrel.Eq{
			ColumnID: filter.IDs,
		})
	}

	if filter.Name != nil {
		builder = builder.Where(squirrel.Like{
			ColumnName: "%" + *filter.Name + "%",
		})
	}

	if filter.Emails != nil {
		builder = builder.Where(squirrel.Eq{
			ColumnEmail: filter.Emails,
		})
	}

	if filter.IsAdmin != nil {
		builder = builder.Where(squirrel.Eq{
			ColumnIsAdmin: *filter.IsAdmin,
		})
	}

	if filter.WithDeleted == nil || *filter.WithDeleted == false {
		builder = builder.Where(squirrel.Eq{
			ColumnDeleted: false,
		})
	}

	if filter.Limit != nil {
		builder = builder.Limit(uint64(*filter.Limit))
	}

	if filter.Offset != nil {
		builder = builder.Offset(uint64(*filter.Offset))
	}

	if filter.Sort != nil {
		builder = builder.OrderBy(*filter.Sort)
	} else {
		builder = builder.OrderBy(ColumnID + " ASC")
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("to sql: %w", err)
	}

	rows, err := r.client.Query(ctx, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("execute query search users: %w", err)
	}
	defer rows.Close()

	users := &Users{}
	users.Result, err = pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[User])
	if err != nil {
		return nil, fmt.Errorf("collect user: %w", err)
	}

	return users, nil
}
