package auth

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"boilerplate/internal/model"
	gin_pkg "boilerplate/internal/pkg/gin"
	jwt_pkg "boilerplate/internal/pkg/jwt"
	"boilerplate/internal/pkg/metadata"
)

var errUnauthorized = errors.New("требуется авторизация")

func (s *service) Validate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accessToken, exists := gin_pkg.GetAccessToken(ctx)
		if exists {
			claims, err := validateToken(accessToken, s.config)
			if err != nil {
				gin_pkg.RenderResponse(ctx, http.StatusUnauthorized, errUnauthorized)
				ctx.Abort()
				return
			}

			userID, exists := GetUserID(claims)
			if !exists {
				gin_pkg.RenderResponse(ctx, http.StatusUnauthorized, errUnauthorized)
				ctx.Abort()
				return
			}

			ctx.Set(metadata.KeyUserID, userID)

			userName, exists := GetUserName(claims)
			if exists {
				ctx.Set(metadata.KeyUserName, userName)
			}

			return
		}

		refreshToken, exists := gin_pkg.GetRefreshToken(ctx)
		if !exists {
			gin_pkg.RenderResponse(ctx, http.StatusUnauthorized, errUnauthorized)
			ctx.Abort()
			return
		}

		claims, err := validateToken(refreshToken, s.config)
		if err != nil {
			gin_pkg.RenderResponse(ctx, http.StatusUnauthorized, errUnauthorized)
			ctx.Abort()
			return
		}

		userID, exists := GetUserID(claims)
		if !exists {
			gin_pkg.RenderResponse(ctx, http.StatusUnauthorized, errUnauthorized)
			ctx.Abort()
			return
		}

		user, err := s.usersService.Get(ctx, userID)
		if err != nil {
			gin_pkg.RenderResponse(ctx, http.StatusUnauthorized, errUnauthorized)
			ctx.Abort()
			return
		}

		if user.Deleted {
			gin_pkg.RenderResponse(ctx, http.StatusUnauthorized, errUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set(metadata.KeyUserID, user.ID)
		ctx.Set(metadata.KeyUserName, user.Name)

		newAccessToken, err := jwt_pkg.GenerateAccessToken(user.ID, user.Name, s.config)
		if err != nil {
			gin_pkg.RenderResponse(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		newRefreshToken, err := jwt_pkg.GenerateRefreshToken(user.ID, user.Name, s.config)
		if err != nil {
			gin_pkg.RenderResponse(ctx, http.StatusInternalServerError, err)
			ctx.Abort()
			return
		}

		gin_pkg.SetAccessToken(ctx, newAccessToken, s.config.AccessTokenTTL)
		gin_pkg.SetRefreshToken(ctx, newRefreshToken, s.config.RefreshTokenTTL)
	}
}

func validateToken(token string, config *model.ConfigAPI) (jwt.MapClaims, error) {
	parsedToken, claims, err := jwt_pkg.ParseToken(token, config)
	if err != nil {
		return nil, errors.New("недействительный токен")
	}

	if !parsedToken.Valid {
		return nil, errors.New("недействительный токен")
	}

	return claims, nil
}

func GetUserID(claims jwt.MapClaims) (int, bool) {
	userID, ok := claims[jwt_pkg.KeyUserID].(float64)
	if !ok {
		return 0, false
	}
	return int(userID), true
}

func GetUserName(claims jwt.MapClaims) (string, bool) {
	userName, ok := claims[jwt_pkg.KeyUserName].(string)
	if !ok {
		return "", false
	}
	return userName, true
}
