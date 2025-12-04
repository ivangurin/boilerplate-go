package auth_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	errors_pkg "boilerplate/internal/pkg/errors"
	suite_factory "boilerplate/internal/pkg/suite/factory"
	suite_provider "boilerplate/internal/pkg/suite/provider"
	"boilerplate/internal/pkg/utils"
	"boilerplate/internal/services/auth"
	"boilerplate/internal/services/users"
)

func TestValidateNoTokens(t *testing.T) {
	sp, cleaner := suite_provider.NewProvider()
	t.Cleanup(cleaner)

	res, err := sp.GetAuthService().Validate(sp.Context(), &auth.AuthValidateRequest{})
	require.Error(t, err)
	require.True(t, errors_pkg.IsErrUnauthorized(err))
	require.Nil(t, res)
}

func TestValidateWrongTokens(t *testing.T) {
	sp, cleaner := suite_provider.NewProvider()
	t.Cleanup(cleaner)

	res, err := sp.GetAuthService().Validate(sp.Context(), &auth.AuthValidateRequest{
		AccessToken: utils.Ptr(gofakeit.UUID()),
	})
	require.Error(t, err)
	require.True(t, errors_pkg.IsErrUnauthorized(err))
	require.Nil(t, res)

	res, err = sp.GetAuthService().Validate(sp.Context(), &auth.AuthValidateRequest{
		RefreshToken: utils.Ptr(gofakeit.UUID()),
	})
	require.Error(t, err)
	require.True(t, errors_pkg.IsErrUnauthorized(err))
	require.Nil(t, res)
}

func TestValidateCorrectTokens(t *testing.T) {
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

	res, err := sp.GetAuthService().Validate(sp.Context(), &auth.AuthValidateRequest{
		AccessToken: utils.Ptr(loginRes.AccessToken),
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	res, err = sp.GetAuthService().Validate(sp.Context(), &auth.AuthValidateRequest{
		RefreshToken: utils.Ptr(loginRes.RefreshToken),
	})
	require.NoError(t, err)
	require.NotNil(t, res)

	err = sp.GetUserService().Delete(sp.Context(), createdUser.ID)
	require.NoError(t, err)

	res, err = sp.GetAuthService().Validate(sp.Context(), &auth.AuthValidateRequest{
		AccessToken: utils.Ptr(loginRes.AccessToken),
	})
	require.Error(t, err)
	require.True(t, errors_pkg.IsErrUnauthorized(err))
	require.Nil(t, res)

	res, err = sp.GetAuthService().Validate(sp.Context(), &auth.AuthValidateRequest{
		RefreshToken: utils.Ptr(loginRes.RefreshToken),
	})
	require.Error(t, err)
	require.True(t, errors_pkg.IsErrUnauthorized(err))
	require.Nil(t, res)
}
