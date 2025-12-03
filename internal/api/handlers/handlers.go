package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swagger_files "github.com/swaggo/files"
	gin_swagger "github.com/swaggo/gin-swagger"

	"boilerplate/internal/api/handlers/auth"
	"boilerplate/internal/api/handlers/auth_public"
	"boilerplate/internal/api/handlers/users"
	"boilerplate/internal/api/middleware"
	"boilerplate/internal/api/swagger"
	logger_pkg "boilerplate/internal/pkg/logger"
	"boilerplate/internal/service_provider"
)

func NewHandler(
	logger logger_pkg.Logger,
	sp *service_provider.Provider,
) http.Handler {
	mw := middleware.NewMiddleware(
		logger,
		sp.GetAuthService(),
	)

	router := gin.New()
	router.Use(mw.Logger())

	// Public routes
	public := router.Group("/api")

	authHandler := auth_public.NewHandler(
		sp.GetAuthService(),
	)

	authHandler.Mount(public.Group("/auth"))

	// Protected routes
	protected := router.Group("/api")
	protected.Use(mw.Auth())

	authHandler = auth.NewHandler(
		sp.GetAuthService(),
	)

	authHandler.Mount(protected.Group("/auth"))

	usersHandler := users.NewHandler(
		sp.GetUsersService(),
	)

	usersHandler.Mount(protected.Group("/users"))

	swagger.SwaggerInfo.BasePath = "/api"
	router.GET("/swagger/*any", gin_swagger.WrapHandler(swagger_files.Handler))

	return router
}
