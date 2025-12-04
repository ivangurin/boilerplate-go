package users_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	suite_factory "boilerplate/internal/pkg/suite/factory"
	suite_provider "boilerplate/internal/pkg/suite/provider"
)

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	sp, cleanup := suite_provider.NewProvider()
	defer cleanup()

	user := suite_factory.NewUserFactory().Build()
	userID, err := sp.GetRepo().Users().Create(sp.Context(), user)
	require.NoError(t, err)

	err = sp.GetUserService().Delete(sp.Context(), userID)
	require.NoError(t, err)

	deletedUser, err := sp.GetUserService().Get(sp.Context(), userID)
	require.NoError(t, err)
	require.NotNil(t, deletedUser)
	require.True(t, deletedUser.Deleted)
	require.NotEmpty(t, deletedUser.DeletedAt)
}
