package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	gin_pkg "boilerplate/internal/pkg/gin"
)

// Delete user
//
//	@Summary		Delete user
//	@Description	Delete user
//	@Tags			UsersAPI
//	@Accept			json
//	@Produce		json
//	@Success		204
//	@Failure		400		{object}	model.HandlerError
//	@Failure		404		{object}	model.HandlerError
//	@Failure		500		{object}	model.HandlerError
//	@Param			userid	path		int	true	"user id"
//	@Router			/users/{userid} [delete]
func (h *handler) Delete(ctx *gin.Context) {
	req := &userRequest{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		gin_pkg.RenderResponse(ctx, http.StatusBadRequest, err)
		return
	}

	err := h.usersService.Delete(ctx, req.UserID)
	if err != nil {
		gin_pkg.RenderErrorResponse(ctx, err)
		return
	}

	gin_pkg.RenderResponse(ctx, http.StatusNoContent, nil)
}
