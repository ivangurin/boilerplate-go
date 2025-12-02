package users

import (
	"boilerplate/internal/repository"
	"context"
)

type Service interface {
	Create(ctx context.Context, req *UserCreateRequest) (*User, error)
	// Get(ctx context.Context, id int) (*User, error)
	// Update(ctx context.Context, id int, req *UserUpdateRequest) (*User, error)
	// Delete(ctx context.Context, id int) error
	// Search(ctx context.Context, req *UserSearchRequest) (*UserSearchResponse, error)
}

type service struct {
	repo repository.Repo
}

func NewService(
	repo repository.Repo,
) Service {
	return &service{
		repo: repo,
	}
}
