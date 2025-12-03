package users_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	suite_factory "boilerplate/internal/pkg/suite/factory"
	suite_provider "boilerplate/internal/pkg/suite/provider"
	"boilerplate/internal/services/users"
)

func TestGetUser(t *testing.T) {
	t.Parallel()

	sp, cleanup := suite_provider.NewProvider()
	defer cleanup()

	user := suite_factory.NewUserFactory().Build()

	createdUser, err := sp.GetUserService().Create(sp.Context(), &users.UserCreateRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	require.NoError(t, err)

	createdUser, err = sp.GetUserService().Get(sp.Context(), createdUser.ID)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.Equal(t, user.Name, createdUser.Name)
	require.Equal(t, user.Email, createdUser.Email)
	require.Equal(t, user.Password, createdUser.Password)
	require.False(t, createdUser.IsAdmin)
	require.False(t, createdUser.Deleted)
	require.NotEmpty(t, createdUser.CreatedAt)
	require.NotEmpty(t, createdUser.UpdatedAt)
}

func TestGetUserNotFound(t *testing.T) {
	t.Parallel()

	sp, cleanup := suite_provider.NewProvider()
	defer cleanup()

	_, err := sp.GetUserService().Get(sp.Context(), -1)
	require.Error(t, err)
}
