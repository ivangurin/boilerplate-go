package middleware

import (
	"github.com/gin-gonic/gin"

	gin_pkg "boilerplate/internal/pkg/gin"
	"boilerplate/internal/services/auth"
)

func (m *middleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := &auth.AuthValidateRequest{}
		accessToken, exists := gin_pkg.GetAccessToken(ctx)
		if exists {
			req.AccessToken = &accessToken
		}
		refreshToken, exists := gin_pkg.GetRefreshToken(ctx)
		if exists {
			req.RefreshToken = &refreshToken
		}

		resp, err := m.authService.Validate(ctx, req)
		if err != nil {
			gin_pkg.RenderErrorResponse(ctx, err)
			ctx.Abort()
			return
		}

		if resp.UserID != nil {
			gin_pkg.SetUserID(ctx, *resp.UserID)
		}

		if resp.UserName != nil {
			gin_pkg.SetUserName(ctx, *resp.UserName)
		}

		if resp.AccessToken != nil {
			gin_pkg.SetAccessToken(ctx, *resp.AccessToken, m.authService.GetConfig().AccessTokenTTL)
		}

		if resp.RefreshToken != nil {
			gin_pkg.SetRefreshToken(ctx, *resp.RefreshToken, m.authService.GetConfig().RefreshTokenTTL)
		}
	}
}
