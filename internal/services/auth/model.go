package auth

import "boilerplate/internal/services/users"

type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	User         *users.User `json:"user"`
}

type AuthRefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthValidateRequest struct {
	AccessToken  *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
}

type AuthValidateResponse struct {
	UserID       *int    `json:"user_id"`
	UserName     *string `json:"user_name"`
	AccessToken  *string `json:"access_token"`
	RefreshToken *string `json:"refresh_token"`
}
