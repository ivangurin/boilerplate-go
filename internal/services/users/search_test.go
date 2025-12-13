package users_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	suite_factory "boilerplate/internal/pkg/suite/factory"
	suite_provider "boilerplate/internal/pkg/suite/provider"
	"boilerplate/internal/pkg/utils"
	users_service "boilerplate/internal/services/users"
)

func TestSearchUsers(t *testing.T) {
	t.Parallel()

	sp, cleanup := suite_provider.NewProvider()
	defer cleanup()

	users := suite_factory.NewUserFactory().Builds(4)
	users[0].IsAdmin = true
	for _, user := range users {
		err := sp.GetRepo().Users().Create(sp.Context(), user)
		require.NoError(t, err)
	}

	err := sp.GetRepo().Users().Delete(sp.Context(), users[3].ID)
	require.NoError(t, err)

	type testCase struct {
		Name     string
		Request  *users_service.UserSearchRequest
		Expected int
	}

	testCases := []testCase{
		{
			Name:     "empty filter",
			Request:  &users_service.UserSearchRequest{},
			Expected: 3,
		},
		{
			Name: "wrong name",
			Request: &users_service.UserSearchRequest{
				Filter: users_service.UserSearchRequestFilter{
					Name: utils.Ptr("wrong"),
				},
			},
			Expected: 0,
		},
		{
			Name: "correct name",
			Request: &users_service.UserSearchRequest{
				Filter: users_service.UserSearchRequestFilter{
					Name: utils.Ptr(users[0].Name),
				},
			},
			Expected: 1,
		},
		{
			Name: "wrong email",
			Request: &users_service.UserSearchRequest{
				Filter: users_service.UserSearchRequestFilter{
					Email: []string{"wrong@example.com"},
				},
			},
			Expected: 0,
		},
		{
			Name: "correct 1 email",
			Request: &users_service.UserSearchRequest{
				Filter: users_service.UserSearchRequestFilter{
					Email: []string{users[0].Email},
				},
			},
			Expected: 1,
		},
		{
			Name: "correct 2 emails",
			Request: &users_service.UserSearchRequest{
				Filter: users_service.UserSearchRequestFilter{
					Email: []string{users[0].Email, users[1].Email},
				},
			},
			Expected: 2,
		},
		{
			Name: "with deleted",
			Request: &users_service.UserSearchRequest{
				Filter: users_service.UserSearchRequestFilter{
					WithDeleted: utils.Ptr(true),
				},
			},
			Expected: 4,
		},
		{
			Name: "limit 1",
			Request: &users_service.UserSearchRequest{
				Limit: utils.Ptr(1),
			},
			Expected: 1,
		},
		{
			Name: "limit 2",
			Request: &users_service.UserSearchRequest{
				Limit: utils.Ptr(2),
			},
			Expected: 2,
		},
		{
			Name: "limit 2 offset 1",
			Request: &users_service.UserSearchRequest{
				Limit:  utils.Ptr(2),
				Offset: utils.Ptr(1),
			},
			Expected: 2,
		},
		{
			Name: "limit 2 offset 2",
			Request: &users_service.UserSearchRequest{
				Limit:  utils.Ptr(2),
				Offset: utils.Ptr(2),
			},
			Expected: 1,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			users, err := sp.GetUserService().Search(sp.Context(), testCase.Request)
			require.NoError(t, err)
			require.Len(t, users.Result, testCase.Expected)
		})
	}
}
