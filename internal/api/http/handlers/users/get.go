package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	gin_pkg "boilerplate/internal/pkg/gin"
)

// Get user
//
//	@Summary		Get user
//	@Description	Get user
//	@Tags			UsersAPI
//	@Accept			json
//	@Produce		json
//	@Success		200		{array}		users.User
//	@Failure		400		{object}	model.HandlerError
//	@Failure		404		{object}	model.HandlerError
//	@Failure		500		{object}	model.HandlerError
//	@Param			userid	path		int	true	"user id"
//	@Router			/users/{userid} [get]
func (h *handler) Get(ctx *gin.Context) {
	req := &userRequest{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		gin_pkg.RenderResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.usersService.Get(ctx, req.UserID)
	if err != nil {
		gin_pkg.RenderErrorResponse(ctx, err)
		return
	}

	gin_pkg.RenderResponse(ctx, http.StatusOK, user)
}
