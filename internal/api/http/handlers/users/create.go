package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	gin_pkg "boilerplate/internal/pkg/gin"
	users_service "boilerplate/internal/services/users"
)

// Create user
//
//	@Summary		Create user
//	@Description	Create user
//	@Tags			UsersAPI
//	@Accept			json
//	@Produce		json
//	@Success		201		{object}	users_service.User
//	@Failure		400		{object}	model.HandlerError
//	@Failure		500		{object}	model.HandlerError
//	@Param			request	body		users_service.UserCreateRequest	true	"user data"
//	@Router			/users [post]
func (h *handler) Create(ctx *gin.Context) {
	req := &users_service.UserCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		gin_pkg.RenderResponse(ctx, http.StatusBadRequest, err)
		return
	}

	user, err := h.usersService.Create(ctx, req)
	if err != nil {
		gin_pkg.RenderErrorResponse(ctx, err)
		return
	}

	gin_pkg.RenderResponse(ctx, http.StatusCreated, user)
}
