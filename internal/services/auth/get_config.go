package auth

import "boilerplate/internal/model"

func (s *service) GetConfig() *model.ConfigAPI {
	return s.config
}
