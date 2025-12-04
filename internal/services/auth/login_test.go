package auth_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"

	errors_pkg "boilerplate/internal/pkg/errors"
	"boilerplate/internal/pkg/jwt"
	suite_factory "boilerplate/internal/pkg/suite/factory"
	suite_provider "boilerplate/internal/pkg/suite/provider"
	"boilerplate/internal/services/auth"
	"boilerplate/internal/services/users"
)

func TestLoginUserNotFound(t *testing.T) {
	sp, cleaner := suite_provider.NewProvider()
	t.Cleanup(cleaner)

	res, err := sp.GetAuthService().Login(sp.Context(), &auth.AuthLoginRequest{
		Email:    gofakeit.Email(),
		Password: gofakeit.Word(),
	})
	require.Error(t, err)
	require.True(t, errors_pkg.IsErrNotFound(err))
	require.Nil(t, res)
}

func TestLoginWrongPassword(t *testing.T) {
	sp, cleaner := suite_provider.NewProvider()
	t.Cleanup(cleaner)

	user := suite_factory.NewUserFactory().WithPassword(gofakeit.Word()).Build()
	_, err := sp.GetUserService().Create(sp.Context(), &users.UserCreateRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	require.NoError(t, err)

	res, err := sp.GetAuthService().Login(sp.Context(), &auth.AuthLoginRequest{
		Email:    user.Email,
		Password: gofakeit.Word(),
	})
	require.Error(t, err)
	require.True(t, errors_pkg.IsErrUnauthorized(err))
	require.Nil(t, res)
}

func TestLoginCorrectPassword(t *testing.T) {
	sp, cleaner := suite_provider.NewProvider()
	t.Cleanup(cleaner)

	user := suite_factory.NewUserFactory().WithPassword(gofakeit.Word()).Build()
	createdUser, err := sp.GetUserService().Create(sp.Context(), &users.UserCreateRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})
	require.NoError(t, err)

	res, err := sp.GetAuthService().Login(sp.Context(), &auth.AuthLoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotEmpty(t, res.AccessToken)
	require.NotEmpty(t, res.RefreshToken)
	require.Equal(t, createdUser.ID, res.User.ID)
	require.Equal(t, createdUser.Name, res.User.Name)
	require.Equal(t, createdUser.Email, res.User.Email)

	claims, err := jwt.ValidateToken(res.AccessToken, sp.GetAuthService().GetConfig())
	require.NoError(t, err)

	userID, exists := jwt.GetUserID(claims)
	require.True(t, exists)
	require.Equal(t, createdUser.ID, userID)

	userName, exists := jwt.GetUserName(claims)
	require.True(t, exists)
	require.Equal(t, createdUser.Name, userName)

	claims, err = jwt.ValidateToken(res.RefreshToken, sp.GetAuthService().GetConfig())
	require.NoError(t, err)

	userID, exists = jwt.GetUserID(claims)
	require.True(t, exists)
	require.Equal(t, createdUser.ID, userID)

	userName, exists = jwt.GetUserName(claims)
	require.True(t, exists)
	require.Equal(t, createdUser.Name, userName)

	err = sp.GetUserService().Delete(sp.Context(), createdUser.ID)
	require.NoError(t, err)

	res, err = sp.GetAuthService().Login(sp.Context(), &auth.AuthLoginRequest{
		Email:    user.Email,
		Password: user.Password,
	})
	require.Error(t, err)
	require.True(t, errors_pkg.IsErrForbidden(err))
	require.Nil(t, res)
}
