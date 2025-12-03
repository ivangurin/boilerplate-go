package users

import (
	"net/http"

	"github.com/gin-gonic/gin"

	gin_pkg "boilerplate/internal/pkg/gin"
	users_service "boilerplate/internal/services/users"
)

// Search users
//
//	@Summary		Search users
//	@Description	Search users
//	@Tags			UsersAPI
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	users_service.User
//	@Failure		400		{object}	model.HandlerError
//	@Failure		404		{object}	model.HandlerError
//	@Failure		500		{object}	model.HandlerError
//	@Param			id		query		[]int	false	"user id"
//	@Param			name	query		string	false	"name"
//	@Param			email	query		string	false	"email"
//	@Param			limit	query		int		false	"limit"
//	@Param			offset	query		int		false	"offset"
//	@Param			sort	query		string	false	"sort"
//	@Router			/users [get]
func (h *handler) Search(ctx *gin.Context) {
	req := &users_service.UserSearchRequest{}
	if err := ctx.ShouldBindQuery(req); err != nil {
		gin_pkg.RenderResponse(ctx, http.StatusBadRequest, err)
		return
	}

	users, err := h.usersService.Search(ctx, req)
	if err != nil {
		gin_pkg.RenderErrorResponse(ctx, err)
		return
	}

	gin_pkg.RenderResponse(ctx, http.StatusOK, users)
}
