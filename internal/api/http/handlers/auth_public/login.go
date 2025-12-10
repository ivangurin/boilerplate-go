package auth_public

import (
	"net/http"

	"github.com/gin-gonic/gin"

	gin_pkg "boilerplate/internal/pkg/gin"
	"boilerplate/internal/services/auth"
)

// Login
//
//	@Summary		Login
//	@Description	Login
//	@Tags			AuthAPI
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	auth.AuthLoginResponse
//	@Failure		400		{object}	model.HandlerError
//	@Failure		401		{object}	model.HandlerError
//	@Failure		404		{object}	model.HandlerError
//	@Failure		500		{object}	model.HandlerError
//	@Param			request	body		auth.AuthLoginRequest	true	"user"
//	@Router			/auth/login [post]
func (h *handler) Login(ctx *gin.Context) {
	req := &auth.AuthLoginRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		gin_pkg.RenderResponse(ctx, http.StatusBadRequest, err)
		return
	}

	res, err := h.authService.Login(ctx, req)
	if err != nil {
		gin_pkg.RenderErrorResponse(ctx, err)
		return
	}

	gin_pkg.SetAccessToken(ctx, res.AccessToken, h.authService.GetConfig().AccessTokenTTL)
	gin_pkg.SetRefreshToken(ctx, res.RefreshToken, h.authService.GetConfig().RefreshTokenTTL)

	gin_pkg.RenderResponse(ctx, http.StatusOK, res)
}
