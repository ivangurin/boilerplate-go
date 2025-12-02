package users

import (
	"context"
	"fmt"
)

func (s *service) Get(ctx context.Context, id int) (*User, error) {
	user, err := s.repo.Users().Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get user: %s", err.Error())
	}

	return toUser(user), nil
}
