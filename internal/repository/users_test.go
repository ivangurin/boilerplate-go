package repository_test

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"

	suite_factory "boilerplate/internal/pkg/suite/factory"
	suite_provider "boilerplate/internal/pkg/suite/provider"
	"boilerplate/internal/pkg/utils"
	"boilerplate/internal/repository"
)

func TestUserCRUD(t *testing.T) {
	t.Parallel()

	sp, cleanup := suite_provider.NewProvider()
	defer cleanup()

	user := suite_factory.NewUserFactory().Build()
	userID, err := sp.GetRepo().Users().Create(sp.Context(), user)
	require.NoError(t, err)
	require.NotZero(t, userID)

	unknownUser, err := sp.GetRepo().Users().Get(sp.Context(), -1)
	require.Error(t, err)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.Nil(t, unknownUser)

	createdUser, err := sp.GetRepo().Users().Get(sp.Context(), userID)
	require.NoError(t, err)
	require.NotNil(t, createdUser)
	require.Equal(t, userID, createdUser.ID)
	require.Equal(t, user.Name, createdUser.Name)
	require.Equal(t, user.Email, createdUser.Email)
	require.Equal(t, user.Password, createdUser.Password)
	require.False(t, createdUser.IsAdmin)
	require.False(t, createdUser.Deleted)
	require.NotEmpty(t, createdUser.CreatedAt)
	require.NotEmpty(t, createdUser.UpdatedAt)

	user = suite_factory.NewUserFactory().Build()
	user.ID = userID
	err = sp.GetRepo().Users().Update(sp.Context(), user)
	require.NoError(t, err)

	updatedUser, err := sp.GetRepo().Users().Get(sp.Context(), userID)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	require.Equal(t, user.Name, updatedUser.Name)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.Password, updatedUser.Password)
	require.False(t, updatedUser.IsAdmin)
	require.False(t, updatedUser.Deleted)
	require.NotEmpty(t, updatedUser.CreatedAt)
	require.NotEmpty(t, updatedUser.UpdatedAt)

	err = sp.GetRepo().Users().Delete(sp.Context(), userID)
	require.NoError(t, err)

	deletedUser, err := sp.GetRepo().Users().Get(sp.Context(), userID)
	require.NoError(t, err)
	require.NotNil(t, deletedUser)
	require.True(t, deletedUser.Deleted)
	require.NotEmpty(t, deletedUser.DeletedAt)
}

func TestUserSearch(t *testing.T) {
	t.Parallel()

	sp, cleanup := suite_provider.NewProvider()
	defer cleanup()

	users := suite_factory.NewUserFactory().Builds(4)
	users[0].IsAdmin = true
	for _, user := range users {
		userID, err := sp.GetRepo().Users().Create(sp.Context(), user)
		require.NoError(t, err)
		user.ID = userID
	}

	err := sp.GetRepo().Users().Delete(sp.Context(), users[3].ID)
	require.NoError(t, err)

	type testCase struct {
		Name     string
		Filter   *repository.UserFilter
		Expected int
	}

	testCases := []testCase{
		{
			Name:     "empty filter",
			Filter:   &repository.UserFilter{},
			Expected: 3,
		},
		{
			Name: "wrong name",
			Filter: &repository.UserFilter{
				Name: utils.Ptr("wrong"),
			},
			Expected: 0,
		},
		{
			Name: "correct name",
			Filter: &repository.UserFilter{
				Name: utils.Ptr(users[0].Name),
			},
			Expected: 1,
		},
		{
			Name: "wrong email",
			Filter: &repository.UserFilter{
				Emails: []string{"wrong@example.com"},
			},
			Expected: 0,
		},
		{
			Name: "correct 1 email",
			Filter: &repository.UserFilter{
				Emails: []string{users[0].Email},
			},
			Expected: 1,
		},
		{
			Name: "correct 2 emails",
			Filter: &repository.UserFilter{
				Emails: []string{users[0].Email, users[1].Email},
			},
			Expected: 2,
		},
		{
			Name: "is admin",
			Filter: &repository.UserFilter{
				IsAdmin: utils.Ptr(true),
			},
			Expected: 1,
		},
		{
			Name: "with deleted",
			Filter: &repository.UserFilter{
				WithDeleted: utils.Ptr(true),
			},
			Expected: 4,
		},
		{
			Name: "limit 1",
			Filter: &repository.UserFilter{
				Limit: utils.Ptr(1),
			},
			Expected: 1,
		},
		{
			Name: "limit 2",
			Filter: &repository.UserFilter{
				Limit: utils.Ptr(2),
			},
			Expected: 2,
		},
		{
			Name: "limit 2 offset 1",
			Filter: &repository.UserFilter{
				Limit:  utils.Ptr(2),
				Offset: utils.Ptr(1),
			},
			Expected: 2,
		},
		{
			Name: "limit 2 offset 2",
			Filter: &repository.UserFilter{
				Limit:  utils.Ptr(2),
				Offset: utils.Ptr(2),
			},
			Expected: 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			users, err := sp.GetRepo().Users().Search(sp.Context(), testCase.Filter)
			require.NoError(t, err)
			require.Len(t, users.Result, testCase.Expected)
		})
	}
}
