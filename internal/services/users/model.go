package users

import (
	"time"

	"boilerplate/internal/model"
	"boilerplate/internal/repository"
)

type User struct {
	ID        int            `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Role      model.UserRole `json:"role"`
	IsAdmin   bool           `json:"is_admin"`
	Password  string         `json:"-"`
	Deleted   bool           `json:"deleted"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt *time.Time     `json:"deleted_at,omitempty"`
}

type UserCreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	ID       int     `uri:"userid"`
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type UserSearchRequest struct {
	Filter UserSearchRequestFilter
	Limit  *int    `form:"limit"`
	Offset *int    `form:"offset"`
	Sort   *string `form:"sort"`
}

type UserSearchRequestFilter struct {
	ID          []int    `form:"id"`
	Name        *string  `form:"name"`
	Email       []string `form:"email"`
	WithDeleted *bool    `form:"with_deleted"`
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
