package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	gin_pkg "boilerplate/internal/pkg/gin"
)

// Me
//
//	@Summary		Me
//	@Description	Hwo am I
//	@Tags			AuthAPI
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	users.User
//	@Failure		400	{object}	model.HandlerError
//	@Failure		401	{object}	model.HandlerError
//	@Failure		404	{object}	model.HandlerError
//	@Failure		500	{object}	model.HandlerError
//	@Router			/auth/me [get]
func (h *handler) Me(ctx *gin.Context) {
	res, err := h.authService.Me(ctx)
	if err != nil {
		gin_pkg.RenderErrorResponse(ctx, err)
		return
	}

	gin_pkg.RenderResponse(ctx, http.StatusOK, res)
}
