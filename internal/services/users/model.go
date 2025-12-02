package users

import (
	"boilerplate/internal/repository"
	"time"
)

type User struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	IsAdmin   bool       `json:"is_admin"`
	Password  string     `json:"-"`
	Deleted   bool       `json:"deleted"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	ID       int
	Name     *string
	Email    *string
	Password *string
}

type UserSearchRequest struct {
	Filter UserSearchRequestFilter
	Limit  *uint64
	Offset *uint64
	Sort   *string
}

type UserSearchRequestFilter struct {
	ID    []int
	Name  *string
	Email []string
}

type UserSearchResponse struct {
	Result []*User `json:"users"`
	Total  int     `json:"total"`
}

func toUser(user *repository.User) *User {
	return &User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		IsAdmin:   user.IsAdmin,
		Password:  user.Password,
		Deleted:   user.Deleted,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		DeletedAt: user.DeletedAt,
	}
}
