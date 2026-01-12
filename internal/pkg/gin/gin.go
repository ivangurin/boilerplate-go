package gin

import (
	"net/http"

	"boilerplate/internal/model"
	"boilerplate/internal/pkg/errors"
	"boilerplate/internal/pkg/metadata"

	"github.com/gin-gonic/gin"
)

const (
	accessToken  = "access_token"
	refreshToken = "refresh_token"
)

func RenderErrorResponse(ctx *gin.Context, err error) {
	if errors.IsErrBadRequest(err) {
		RenderResponse(ctx, http.StatusBadRequest, err)
		return
	}
	if errors.IsErrNotFound(err) {
		RenderResponse(ctx, http.StatusNotFound, err)
		return
	}
	if errors.IsErrForbidden(err) {
		RenderResponse(ctx, http.StatusForbidden, err)
		return
	}
	if errors.IsErrUnauthorized(err) {
		RenderResponse(ctx, http.StatusUnauthorized, err)
		return
	}

	RenderResponse(ctx, http.StatusInternalServerError, err)
}

func RenderResponse(ctx *gin.Context, status int, resp any) {
	switch v := resp.(type) {
	case string:
		ctx.JSON(status, &model.HandlerMessage{
			Message: v,
		})
	case error:
		ctx.JSON(status, &model.HandlerError{
			Error: v.Error(),
		})

	case nil:
		ctx.Status(status)

	default:
		ctx.JSON(status, resp)
	}
}

func setCookie(ctx *gin.Context, name, value string, maxAge int, path string, secure, httpOnly bool) {
	domain := ""
	ctx.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
}

func SetAccessToken(ctx *gin.Context, value string, ttl int) {
	setCookie(ctx, accessToken, value, ttl, "/", false, true)
}

func SetRefreshToken(ctx *gin.Context, value string, ttl int) {
	setCookie(ctx, refreshToken, value, ttl, "/", false, true)
}

func GetAccessToken(ctx *gin.Context) (string, bool) {
	value, err := ctx.Cookie(accessToken)
	if err != nil {
		return "", false
	}
	return value, true
}

func GetRefreshToken(ctx *gin.Context) (string, bool) {
	value, err := ctx.Cookie(refreshToken)
	if err != nil {
		return "", false
	}
	return value, true
}

func SetUserID(ctx *gin.Context, userID int) {
	ctx.Set(metadata.KeyUserID, userID)
}
