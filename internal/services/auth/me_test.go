package auth_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	errors_pkg "boilerplate/internal/pkg/errors"
	gin_pkg "boilerplate/internal/pkg/gin"
	suite_factory "boilerplate/internal/pkg/suite/factory"
	suite_provider "boilerplate/internal/pkg/suite/provider"
)

func TestMe_NoUserID(t *testing.T) {
	sp, cleaner := suite_provider.NewProvider()
	t.Cleanup(cleaner)

	rw := httptest.NewRecorder()
	gCtx, _ := gin.CreateTestContext(rw)

	user, err := sp.GetAuthService().Me(gCtx)
	require.Error(t, err)
	require.Nil(t, user)
	require.True(t, errors_pkg.IsErrUnauthorized(err))
}

func TestMe_CorrectUser(t *testing.T) {
	sp, cleaner := suite_provider.NewProvider()
	t.Cleanup(cleaner)

	user := suite_factory.NewUserFactory().Build()
	err := sp.GetRepo().Users().Create(sp.Context(), user)
	require.NoError(t, err)

	rw := httptest.NewRecorder()
	gCtx, _ := gin.CreateTestContext(rw)
	gin_pkg.SetUserID(gCtx, user.ID)

	res, err := sp.GetAuthService().Me(gCtx)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, user.ID, res.ID)
	require.Equal(t, user.Name, res.Name)
	require.Equal(t, user.Email, res.Email)
}
