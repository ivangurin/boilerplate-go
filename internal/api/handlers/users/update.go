package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	gin_pkg "boilerplate/internal/pkg/gin"
	users_service "boilerplate/internal/services/users"
)

// Update user
//
//	@Summary		Update user
//	@Description	Update user
//	@Tags			UsersAPI
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	users_service.User
//	@Failure		400		{object}	model.HandlerError
//	@Failure		404		{object}	model.HandlerError
//	@Failure		500		{object}	model.HandlerError
//	@Param			request	body		users_service.UserUpdateRequest	true	"user data"
//	@Router			/users [patch]
func (h *handler) Update(ctx *gin.Context) {
	req := &users_service.UserUpdateRequest{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		gin_pkg.RenderResponse(ctx, http.StatusBadRequest, err)
		return
	}
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		gin_pkg.RenderResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.usersService.Update(ctx, req)
	if err != nil {
		gin_pkg.RenderErrorResponse(ctx, err)
		return
	}

	gin_pkg.RenderResponse(ctx, http.StatusOK, user)
}
