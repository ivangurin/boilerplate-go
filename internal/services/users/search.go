package users

import (
	"boilerplate/internal/repository"
	"context"
	"fmt"
)

func (s *service) Search(ctx context.Context, req *UserSearchRequest) (*UserSearchResponse, error) {
	filter := &repository.UserFilter{
		ID:     req.Filter.ID,
		Email:  req.Filter.Email,
		Name:   req.Filter.Name,
		Limit:  req.Limit,
		Offset: req.Offset,
		Sort:   req.Sort,
	}

	users, err := s.repo.Users().Search(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get user: %s", err.Error())
	}

	resp := &UserSearchResponse{
		Result: make([]*User, 0, len(users.Result)),
	}

	for _, u := range users.Result {
		resp.Result = append(resp.Result, toUser(u))
	}

	return resp, nil
}
