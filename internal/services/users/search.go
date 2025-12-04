package users

import (
	"boilerplate/internal/repository"
	"context"
	"fmt"
)

func (s *service) Search(ctx context.Context, req *UserSearchRequest) (*UserSearchResponse, error) {
	filter := &repository.UserFilter{
		IDs:         req.Filter.ID,
		Emails:      req.Filter.Email,
		Name:        req.Filter.Name,
		WithDeleted: req.Filter.WithDeleted,
		Limit:       req.Limit,
		Offset:      req.Offset,
		Sort:        req.Sort,
	}

	users, err := s.repo.Users().Search(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	resp := &UserSearchResponse{
		Result: make([]*User, 0, len(users.Result)),
	}

	for _, u := range users.Result {
		resp.Result = append(resp.Result, toUser(u))
	}

	return resp, nil
}
