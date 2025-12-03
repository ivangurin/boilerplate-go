package users

import (
	"boilerplate/internal/model"
	"boilerplate/internal/services/users"

	"github.com/gin-gonic/gin"
)

type handler struct {
	usersService users.Service
}

func NewHandler(
	usersService users.Service,
) model.Handler {
	return &handler{
		usersService: usersService,
	}
}

func (h *handler) Mount(router *gin.RouterGroup) {
	router.POST("/", h.Create)
	router.GET("/:userid", h.Get)
	router.PATCH("/:userid", h.Update)
	router.DELETE("/:userid", h.Delete)
	router.GET("/", h.Search)
}
