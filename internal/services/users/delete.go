package users

import (
	"context"
	"fmt"
)

func (s *service) Delete(ctx context.Context, id int) error {
	_, err := s.repo.Users().Get(ctx, id)
	if err != nil {
		return fmt.Errorf("get user: %s", err.Error())
	}

	err = s.repo.Users().Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("delete user: %s", err.Error())
	}

	return nil
}
