package auth_public

import (
	"boilerplate/internal/model"
	"boilerplate/internal/services/auth"

	"github.com/gin-gonic/gin"
)

type handler struct {
	authService auth.Service
}

func NewHandler(
	authService auth.Service,
) model.Handler {
	return &handler{
		authService: authService,
	}
}

func (h *handler) Mount(router *gin.RouterGroup) {
	router.POST("/login", h.Login)
}
