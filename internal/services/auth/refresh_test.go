package auth_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	errors_pkg "boilerplate/internal/pkg/errors"
	suite_factory "boilerplate/internal/pkg/suite/factory"
	suite_provider "boilerplate/internal/pkg/suite/provider"
	"boilerplate/internal/services/auth"
	"boilerplate/internal/services/users"
)

func TestRefreshNoToken(t *testing.T) {
	sp, cleaner := suite_provider.NewProvider()
	t.Cleanup(cleaner)

	res, err := sp.GetAuthService().Refresh(sp.Context(), &auth.AuthRefreshRequest{})
	require.Error(t, err)
	require.True(t, errors_pkg.IsErrBadRequest(err))
	require.Nil(t, res)
}

func TestRefreshNotCorrectToken(t *testing.T) {
	sp, cleaner := suite_provider.NewProvider()
	t.Cleanup(cleaner)

	res, err := sp.GetAuthService().Refresh(sp.Context(), &auth.AuthRefreshRequest{
		RefreshToken: gofakeit.UUID(),
	})
	require.Error(t, err)
	require.True(t, errors_pkg.IsErrUnauthorized(err))
	require.Nil(t, res)
}

func TestRefreshCorrectToken(t *testing.T) {
	sp, cleaner := suite_provider.NewProvider()
	t.Cleanup(cleaner)

	user := suite_factory.NewUserFactory().WithPassword(gofakeit.Word()).Build()
	createdUser, err := sp.GetUserService().Create(sp.Context(), &users.UserCreateRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	require.NoError(t, err)

	loginRes, err := sp.GetAuthService().Login(sp.Context(), &auth.AuthLoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})
	require.NoError(t, err)
	require.NotNil(t, loginRes)

	res, err := sp.GetAuthService().Refresh(sp.Context(), &auth.AuthRefreshRequest{
		RefreshToken: loginRes.RefreshToken,
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.AccessToken)
	require.NotEmpty(t, res.RefreshToken)

	err = sp.GetUserService().Delete(sp.Context(), createdUser.ID)
	require.NoError(t, err)

	res, err = sp.GetAuthService().Refresh(sp.Context(), &auth.AuthRefreshRequest{
		RefreshToken: loginRes.RefreshToken,
	})
	require.Error(t, err)
	require.True(t, errors_pkg.IsErrForbidden(err))
	require.Nil(t, res)
}
