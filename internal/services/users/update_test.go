package users_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	suite_factory "boilerplate/internal/pkg/suite/factory"
	suite_provider "boilerplate/internal/pkg/suite/provider"
	"boilerplate/internal/services/users"
)

func TestUpdateUser(t *testing.T) {
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

	user = suite_factory.NewUserFactory().Build()

	updatedUser, err := sp.GetUserService().Update(sp.Context(), &users.UserUpdateRequest{
		ID:       createdUser.ID,
		Name:     &user.Name,
		Email:    &user.Email,
		Password: &user.Password,
	})
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	require.Equal(t, user.Name, updatedUser.Name)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.Password, updatedUser.Password)
	require.False(t, updatedUser.IsAdmin)
	require.False(t, updatedUser.Deleted)
	require.NotEmpty(t, updatedUser.CreatedAt)
	require.NotEmpty(t, updatedUser.UpdatedAt)
}
