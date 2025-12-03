package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	gin_pkg "boilerplate/internal/pkg/gin"
)

// Logout
//
//	@Summary		Logout
//	@Description	Logout
//	@Tags			AuthAPI
//	@Accept			json
//	@Produce		json
//	@Success		200
//	@Router			/auth/logout [post]
func (h *handler) Logout(ctx *gin.Context) {
	gin_pkg.SetAccessToken(ctx, "", -1)
	gin_pkg.SetRefreshToken(ctx, "", -1)
	gin_pkg.RenderResponse(ctx, http.StatusOK, nil)
}
